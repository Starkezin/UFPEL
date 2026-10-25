package ast;

import java.util.ArrayList;

public class CIf extends Comando{
	public int linha;
	public Exp exp;
	public ArrayList<Comando> bloco;
	
	public CIf(int linha,Exp exp, ArrayList<Comando> bloco)
	{
	  this.linha = linha;
	  this.exp = exp;
	  this.bloco = bloco;
	} 

	@Override
    public String geraCodigo() {
        StringBuilder sb = new StringBuilder();
        
        // Gera apenas o IF
        sb.append("if (").append(exp.geraCodigo()).append(") {\n");
        
        if (bloco != null) {
            for (Comando c : bloco) {
                sb.append(c.geraCodigo());
            }
        }
        
        sb.append("}\n"); 
        // REMOVIDO: Qualquer geração de 'else' automática
        
        return sb.toString();
    }

}
