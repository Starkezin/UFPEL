# FAST FRUIT PROTOCOL

## Descrição
Jogo multiplayer que funciona em rede local. É necessário configurar o servidor e conectar dois dispositivos para jogar.

## Passo 1: Configurar o cliente
1. Abra o arquivo `client.kt`.
2. Altere o valor de `host` para o IPv4 da rede à qual os dispositivos estão conectados.

## Passo 2: Compilar e executar o servidor
1. Compile o servidor com o Kotlin Compiler:

```bash
kotlinc server.kt -include-runtime -d server.jar
```
2. Execute o servidor:

```bash
java -jar server.jar
```

## Passo 3: Jogar
1. Instale o APK do jogo em dois dispositivos diferentes (android).
2. Abra o jogo em ambos os aparelhos.
3. Clique em "Iniciar Jogo" para que os dispositivos se pareiem automaticamente.
