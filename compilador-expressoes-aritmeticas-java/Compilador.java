import java.io.*;

class Compilador{

	public static void main(String[]args)
	{	
		ArvoreSintatica arv=null;
	
		try{

			AnaliseLexica al = new AnaliseLexica(args[0]);
			Parser as = new Parser(al);
		
			arv = as.parseProg();
		
			
			CodeGen backend = new CodeGen();
			String codigo = backend.geraCodigo(arv);

            String outputFile = "output.txt";
            
            try (FileWriter writer = new FileWriter(outputFile)) {
                writer.write(codigo);
                System.out.println("Compilado em: " + outputFile);
            } catch (IOException e) {
                System.out.println("Erro ao salvar arquivo: " + e.getMessage());
            }


		}catch(Exception e)
		{			
			System.out.println("Erro de compilação:\n" + e);
		}
	}
}
