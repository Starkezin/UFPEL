package com.example.fastfruit

import android.app.AlertDialog
import android.app.Dialog
import android.content.Intent
import android.os.Bundle
import androidx.core.content.ContentProviderCompat.requireContext
import androidx.fragment.app.DialogFragment

class PlacarDialogFragment : DialogFragment() {

    companion object {
        fun fromString(placar: String): PlacarDialogFragment {
            val fragment = PlacarDialogFragment()
            val args = Bundle()
            args.putString("placar", placar)
            fragment.arguments = args
            return fragment
        }
    }

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val placar = arguments?.getString("placar") ?: "Sem dados"
        return AlertDialog.Builder(requireContext())
            .setTitle("Placar Final")
            .setMessage(placar)
            .setPositiveButton("OK") { dialog, _ ->
                dialog.dismiss()
                val intent = Intent(requireContext(), HomeScreen::class.java)
                intent.addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP or Intent.FLAG_ACTIVITY_NEW_TASK)
                startActivity(intent)
                requireActivity().finish()
            }
            .create()
    }
}