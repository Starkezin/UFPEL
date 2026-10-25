package ast;

import java.util.ArrayList;

public class EChamadaFun extends Exp{
	public String fun;
	public ArrayList<Exp> args;
	
	public EChamadaFun(String fun, ArrayList<Exp> args)
	{
	  this.fun = fun;
	  this.args = args;
	} 

	@Override
	public String geraCodigo() {
		String codigo = fun + "(";
		for (int i = 0; i < args.size(); i++) {
			codigo += args.get(i).geraCodigo();
			if (i < args.size() - 1) {
				codigo += ", ";
			}
		}
		return codigo + ")";
	}

}
