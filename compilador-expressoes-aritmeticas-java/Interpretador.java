class Interpretador {

    public int interpretador(ArvoreSintatica tree){
        if (tree instanceof Num) {
            return ((Num) tree).num;
        }

        if (tree instanceof Operador) {
            Operador op = (Operador) tree;
            int left = interpretador(op.arg1);
            int right = interpretador(op.arg2);

            if (op instanceof Soma) return left + right;
            if (op instanceof Sub) return left - right;
            if (op instanceof Mult) return left * right;
            if (op instanceof Div) return left / right;

            throw new IllegalArgumentException("Operador desconhecido: " + op.getClass().getName());

        }
        throw new IllegalArgumentException("Tipo de nó desconhecido: " + tree.getClass().getName());
    }

    public static void main(String[]args)
	{	
		ArvoreSintatica arv=null;
	
		try{

			AnaliseLexica al = new AnaliseLexica(args[0]);
			Parser as = new Parser(al);
		
			arv = as.parseProg();

			Interpretador interp = new Interpretador();
			int resultado = interp.interpretador(arv);
			System.out.println("Resultado: " + resultado);

		}catch(Exception e)
		{			
			System.out.println("Erro de compilação:\n" + e);
		}
	}
 
}