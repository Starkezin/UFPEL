package ast;

public class VarDecl{
   public String type;
   public String var;
   
   public VarDecl(String type,String var){
   	this.type = type;
   	this.var = var;
   }

   public String geraCodigo() {
        // 1. Tradução de Tipos
        // O C usa "int", "float", "char". Se sua linguagem usa outros nomes,
        // você precisa converter aqui.
        String tipoC = type; // Assume que já é compatível por padrão
        
        if (type.equals("Float")){
            tipoC = "float";
        } else if (type.equals("Bool")){
            tipoC = "int";
        } else if (type.equals("Void")) {
            tipoC = "void"; // Ou char id[100], dependendo da sua implementação
        }

        // 2. Retorna a declaração: tipo nome;
        // Exemplo: int contador;
        return tipoC + " " + var + ";\n";
    }
}
