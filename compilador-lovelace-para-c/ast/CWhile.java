package ast;

import java.util.ArrayList;

public class CWhile extends Comando{
	public int linha;
	public Exp exp;
	public ArrayList<Comando> bloco;
	
	public CWhile(int linha,Exp exp, ArrayList<Comando> bloco)
	{
	  this.linha = linha;
	  this.exp = exp;
	  this.bloco = bloco;
	} 

	@Override
	public String geraCodigo(){
		String s = "while (" + exp.geraCodigo() + "){\n";
		for (Comando c: bloco){
			s += " " + c.geraCodigo() + "\n";
		}
		s += "}";
		return s;
	}

}
