package ast;

public class ParamFormalFun{
   public String type;
   public String var;
   
   public ParamFormalFun(String type,String var){
   	this.type = type;
   	this.var = var;
   }

   public String geraCodigo(){
   	String tipoC = type;

      if(type.equals("Float")) tipoC = "float";
      else if(type.equals("Bool")) tipoC = "int";
      else if(type.equals("Void")) tipoC = "void";
      return tipoC + " " + var;
   }
}
