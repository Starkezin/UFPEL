package com.example.fastfruit

import java.io.BufferedReader
import java.io.InputStreamReader
import java.net.ServerSocket
import kotlin.concurrent.thread
import java.io.PrintWriter
import java.io.ObjectOutputStream
import kotlin.random.Random

fun main() {
    val totalFrutas = 22
    val qtdSelecionadas = 16
    var frutasSorteadas = (0 until totalFrutas).shuffled().take(qtdSelecionadas)
    var frutasDestaque = frutasSorteadas.shuffled().take(11)

    var rodada = 0

    val acertosPorPlayer = mutableMapOf<Int,Int>()

    val serverSocket = ServerSocket(12345)
    println("Servidor rodando na porta 12345...")

    var numPlayersConnected = 0
    val clients = mutableMapOf<Int, ObjectOutputStream>()
    var nextId = 1

    while (true) {
        val clientSocket = serverSocket.accept()
        val id = nextId++
        synchronized(acertosPorPlayer){
            acertosPorPlayer[id] = 0
        }

        println("Novo cliente conectado!, $id")

        thread {
            val oos = ObjectOutputStream(clientSocket.getOutputStream())
            val reader = BufferedReader(InputStreamReader(clientSocket.getInputStream()))

            synchronized(clients) { clients[id] = oos }

            try {
                while (true) {
                    val mensagem = reader.readLine() ?: break

                    when (mensagem) {
                        "READY" -> {
                            println("READY CHEGOU!")
                            numPlayersConnected++
                            println("Jogadores prontos: $numPlayersConnected")

                            if (numPlayersConnected == 2) {
                                frutasSorteadas = (0 until totalFrutas).shuffled().take(qtdSelecionadas)
                                println("FRUTAS: $frutasSorteadas")
                                frutasDestaque = frutasSorteadas.shuffled().take(11)
                                println("FRUTAS DESTAQUE: $frutasDestaque")

                                rodada = 0 // resetar rodada para nova partida
                                synchronized(acertosPorPlayer) {
                                    acertosPorPlayer.keys.forEach { acertosPorPlayer[it] = 0 }
                                }
                                println("Todos prontos, enviando frutas...")
                                numPlayersConnected = 0
                                synchronized(clients) {
                                    clients.values.forEach { out ->
                                        out.writeObject(frutasSorteadas)
                                        out.writeObject(frutasDestaque)
                                        out.flush()
                                    }
                                }
                            }
                        }
                        "STARTGAME_OK" -> {
                            println("STARTGAME_OK CHEGOU!")
                        }
                        "SENDFRUIT" -> {
                            val msg = reader.readLine().toInt()
                            println("$msg")

                            println("SENDFRUIT FROM JOGADOR, $id")

                            if(msg == frutasDestaque[rodada]){

                                synchronized(acertosPorPlayer){
                                    acertosPorPlayer[id] = acertosPorPlayer.getOrDefault(id, 0) + 1
                                }

                                println("$rodada")
                                synchronized(clients) {
                                    clients.values.forEach { out ->
                                        out.writeObject("RIGHTANSWER")
                                        out.flush()
                                    }
                                }

                                println("Placar:")
                                acertosPorPlayer.forEach { (id, acertos) ->
                                    println("Jogador $id -> $acertos acertos")
                                }
                                rodada++
                                if(rodada >= 3){
                                    println("aqui")
                                    val placar = acertosPorPlayer.entries.joinToString("\n") { (id, pontos) ->
                                        "Player $id: $pontos"
                                    }
                                    synchronized(clients){
                                        clients.values.forEach { out ->
                                            out.writeObject("GAMEFINISH")
                                            out.writeObject(placar)
                                            out.flush()
                                        }
                                    }
                                }

                            }
                            else {
                                println("Jogador $id Errou!!!!")
                            }
                        }
                    }
                }
            } catch (e: Exception) {
                println("Erro com cliente: ${e.message}")
            } finally {
                synchronized(clients) { clients.remove(id) }
                clientSocket.close()
                println("Cliente $id desconectado, ainda conectados: ${clients.size}")
            }
        }
    }
}
