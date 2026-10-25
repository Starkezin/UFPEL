package com.example.fastfruit

import androidx.appcompat.app.AppCompatActivity
import java.io.ObjectInputStream
import java.io.PrintWriter
import java.net.Socket
import java.io.BufferedReader
import java.io.InputStreamReader

class Client(private val host: String = "192.168.8.23", private val port: Int = 12345) {

    private lateinit var socketCliente: Socket
    private lateinit var writer: PrintWriter
    private lateinit var reader: BufferedReader
    private lateinit var ois: ObjectInputStream
    var isConnected = false
        private set

    fun connect(onConnected: () -> Unit){
        Thread{
            try {
                socketCliente = Socket(host, port)
                writer = PrintWriter(socketCliente.getOutputStream(), true)
                reader = BufferedReader(InputStreamReader(socketCliente.getInputStream()))
                ois = ObjectInputStream(socketCliente.getInputStream())
                println("Bem vindo ao servidor")
                isConnected = true
                onConnected()
            } catch (e: Exception){
                e.printStackTrace()
            }
        }.start()
    }

    fun sendReady(onGameStart: (List<Int>, List<Int>) -> Unit) {
        if (!isConnected) return
        Thread {
            writer.println("READY")
            val frutasRecebidas = ois.readObject() as ArrayList<Int>
            val frutasDestaque = ois.readObject() as ArrayList<Int>
            onGameStart(frutasRecebidas, frutasDestaque)
        }.start()
    }

    fun sendStartGame_OK(){
        Thread{
            writer.println("STARTGAME_OK")
        }.start()
    }

    fun enviarClickParaServidor(index: Int) {
        Thread {
            writer.println("SENDFRUIT")
            writer.flush()
            writer.println(index)
            writer.flush()
        }.start()
    }

    fun startListener(onMessage: (String) -> Unit, onFinish: (String) -> Unit) {
        Thread {
            while (true) {
                val msg = ois.readObject()
                if (msg is String) {
                    when {
                        msg == "GAMEFINISH" -> {
                            val placar = ois.readObject().toString()
                            onFinish(placar)
                        }
                        else -> {
                            onMessage(msg)
                        }
                    }
                }
            }
        }.start()
    }

}
