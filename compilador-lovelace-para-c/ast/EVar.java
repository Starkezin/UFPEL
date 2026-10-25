package ast;


public class EVar extends Exp{
	public String var;
	
	
	public EVar(String var)
	{
	  this.var = var;
	} 
	
	
	@Override
	public String geraCodigo(){
		return var;
	}
}
