package com.example.fastfruit

import android.content.Intent
import android.os.Bundle
import android.media.MediaPlayer
import android.view.View
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import java.util.ArrayList

val client = Client()
class HomeScreen : AppCompatActivity() {



    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_home_screen)

        val startButton = findViewById<View>(R.id.startButton)

        client.connect {
            runOnUiThread {
                Toast.makeText(this, "Conectado ao servidor!", Toast.LENGTH_SHORT).show()
            }
        }

        startButton.setOnClickListener{
            client.sendReady { frutas, frutasDestaque ->
                runOnUiThread {
                    val intent = Intent(this, GameScreen::class.java)
                    intent.putIntegerArrayListExtra("frutas", ArrayList(frutas))
                    intent.putIntegerArrayListExtra("frutasDestaque", ArrayList(frutasDestaque))
                    startActivity(intent)
                }
            }
        }
    }
}
