package ast;

import java.util.ArrayList;

public class Fun{
	public String nome;
	public ArrayList<ParamFormalFun> params;
	public String retorno;
	public ArrayList<VarDecl> vars;
	public ArrayList<Comando> body;
	
	public Fun(String nome,ArrayList<ParamFormalFun> params, String retorno,ArrayList<VarDecl> vars,ArrayList<Comando> body)
	{
		this.nome = nome;
		this.params = params;
		this.retorno = retorno;
		this.vars = vars;
		this.body = body;
	}

	public String geraCodigo() {
        // 1. Traduz o tipo de retorno da função (Float -> float)
        String tipoRetornoC = retorno;
        if (retorno.equals("Float")) tipoRetornoC = "float";
        else if (retorno.equals("Bool")) tipoRetornoC = "int";
        else if (retorno.equals("Void")) tipoRetornoC = "void";

        StringBuilder sb = new StringBuilder();
        
        // Ex: float soma(
        sb.append(tipoRetornoC).append(" ").append(nome).append("(");
        
        for (int i = 0; i < params.size(); i++) {
            sb.append(params.get(i).geraCodigo()); 
            
            if (i < params.size() - 1) {
                sb.append(", ");
            }
        }
        
        sb.append(") {\n");
        
        // Variáveis locais
        for (VarDecl v : vars) {
            sb.append(v.geraCodigo());
        }
        
        // Corpo da função
        for (Comando c : body) {
            sb.append(c.geraCodigo());
        }
        
        sb.append("}\n\n");
        return sb.toString();
    }
}
