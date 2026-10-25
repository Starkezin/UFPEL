package ast;

import java.util.ArrayList;

public class Main{

	public ArrayList<VarDecl> vars;
	public ArrayList<Comando> coms;
	
	public Main(ArrayList<VarDecl> vars,ArrayList<Comando> coms)
	{
		this.vars = vars;
		this.coms = coms;
	}

	public String geraCodigo() {
        StringBuilder sb = new StringBuilder();
        
        // Assinatura padrão da main em C
        sb.append("int main() {\n");
        
        for (VarDecl v : vars) {
            sb.append(v.geraCodigo());
        }

        for (Comando c : coms) {
            sb.append(c.geraCodigo());
        }
        
        // Boa prática: return 0 no final da main
        sb.append("\nreturn 0;\n");
        sb.append("}\n");
        
        return sb.toString();
    }
}
