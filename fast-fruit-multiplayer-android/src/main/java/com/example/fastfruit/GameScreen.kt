package com.example.fastfruit

import android.os.Bundle
import android.os.CountDownTimer
import android.os.Handler
import android.os.Looper
import android.view.View
import android.widget.ImageView
import android.widget.TextView
import androidx.activity.enableEdgeToEdge
import androidx.appcompat.app.AppCompatActivity
import androidx.core.view.isVisible
import java.time.ZonedDateTime

class GameScreen : AppCompatActivity() {

    private lateinit var imageButtons: List<ImageView>
    private lateinit var frutasEmJogo: MutableList<Int>


    private val frutas = arrayOf(
        0 to R.drawable.apple, 1 to R.drawable.banana, 2 to R.drawable.cherry, 3 to R.drawable.coconut,
        4 to R.drawable.grape, 5 to R.drawable.kiwi, 6 to R.drawable.lemon, 7 to R.drawable.mango,
        8 to R.drawable.orange, 9 to R.drawable.papaya, 10 to R.drawable.pear, 11 to R.drawable.persimmon,
        12 to R.drawable.pineapple, 13 to R.drawable.strawberry, 14 to R.drawable.tomato, 15 to R.drawable.watermelon,
        16 to R.drawable.plum, 17 to R.drawable.sourspop, 18 to R.drawable.peach, 19 to  R.drawable.guava,
        20 to R.drawable.green_apple, 21 to R.drawable.green_grape
    )

    private lateinit var frutasServer: ArrayList<Int>
    private lateinit var frutasDestaqueS: ArrayList<Int>
    private var rodada = 0


    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        frutasServer = intent.getIntegerArrayListExtra("frutas") ?: ArrayList()
        frutasDestaqueS = intent.getIntegerArrayListExtra("frutasDestaque") ?: ArrayList()

        startGame()
    }

    private fun startGame() {
        setContentView(R.layout.activity_game_screen_nivel_three)
        getImagesButton()
        frutasEmJogo = frutasServer.toMutableList()
        println("Frutas: $frutasEmJogo")
        client.sendStartGame_OK()
        definirImagemDestaque(frutasDestaqueS[0])


        posicionarImagens()
        jogar()
    }

    private fun qtdFrutas(): Int {
        return 16
    }

    private fun getImagesButton() {
        imageButtons = (1..qtdFrutas()).mapNotNull { i ->
            val id = resources.getIdentifier("fruit_$i", "id", packageName)
            findViewById(id)
        }
    }


    private fun definirImagemDestaque(imagem: Int) {
        val drawableID = frutas[imagem]
        trocarImagem(findViewById(R.id.fruta_destaque), drawableID.second)
    }

    private fun trocarImagem(imagemAtual: ImageView, idNovaImagem: Int) {
        imagemAtual.setImageResource(idNovaImagem)
    }

    private fun posicionarImagens() {
        getImagesButton()
        imageButtons.forEachIndexed { index, image ->
            val img = frutasEmJogo
            val drawbleID = frutas[img[index]].second
            trocarImagem(image, drawbleID)
        }
    }

    private fun removerImagem(image: ImageView, index: Int) {
        Handler(Looper.getMainLooper()).postDelayed({
            if (index in frutasEmJogo.indices) {
                frutasEmJogo[index] = -1
                image.visibility = View.INVISIBLE
                image.isEnabled = false
            }
        }, 100)
    }


    private fun jogar() {
        val aux = frutasEmJogo.toList()

        client.startListener(
            onMessage = { msg ->
                runOnUiThread {
                    val index = frutasEmJogo.indexOf(frutasDestaqueS[rodada])
                    removerImagem(imageButtons[index], frutasDestaqueS[rodada])

                    val countdownText = findViewById<TextView>(R.id.countdownText)

                    countdownText.visibility = View.VISIBLE
                    if (rodada < 3) {
                        object : CountDownTimer(4000, 1000) {
                            override fun onTick(millisUntilFinished: Long) {
                                val seconds = (millisUntilFinished / 1000).toInt()
                                countdownText.text = if (seconds > 0) "$seconds" else "GO!"
                            }

                            override fun onFinish() {
                                rodada++
                                definirImagemDestaque(frutasDestaqueS[rodada])

                                countdownText.postDelayed({
                                    countdownText.visibility = View.GONE
                                }, 1000)
                            }
                        }.start()
                    }
                }
            },
            onFinish = { placar ->
                runOnUiThread {
                    val dialog = PlacarDialogFragment.fromString(placar)
                    dialog.show(supportFragmentManager, "dialogPlacar")
                }
            }
        )

        imageButtons.forEachIndexed { index, image ->
            image.setOnClickListener {
                val frutaClicada = aux[index]
                client.enviarClickParaServidor(frutaClicada)
            }
        }
    }




}
