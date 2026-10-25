import java.io.*;

class MaquinaPilha {
    private int[] stack;
    private int top;
    private static final int STACK_SIZE = 100;
    

    public MaquinaPilha() {
        stack = new int[STACK_SIZE];
        top = -1;
    }

    private void push(int value){
        if (top < STACK_SIZE -1){
            stack[++top] = value;
        } else {
            throw new StackOverflowError("Pilha cheia");
        }
    }

    private int pop(){
        if (top >= 0){
            return stack[top--];
        } else {
            throw new IllegalStateException("Pilha vazia");
        }
    }

    public void execute(String[] instructions) {
        for (String inst : instructions) {
            String[] parts = inst.split(" ");
            String command = parts[0];

            switch (command) {
                case "PUSH":
                    int value = Integer.parseInt(parts[1]);
                    push(value);
                    break;
                case "SUM":
                    int b = pop();
                    int a = pop();
                    push(a + b);
                    break;
                case "SUB":
                    b = pop();
                    a = pop();
                    push(a - b);
                    break;
                case "MULT":
                    b = pop();
                    a = pop();
                    push(a * b);
                    break;
                case "DIV":
                    b = pop();
                    a = pop();
                    push(a / b);
                    break;
                case "PRINT":
                    System.out.println(pop());
                    break;
                default:
                    throw new IllegalArgumentException("Instrução desconhecida: " + command);
            }
        }
    }

    public int getResult(){
        return stack[top];
    }



    public static void main(String[] args) {
       if (args.length != 1) {
            System.out.println("Uso: java MaquinaPilha <output.txt>");
            return;
        }

        try {
            // Lê o arquivo de entrada
            BufferedReader reader = new BufferedReader(new FileReader(args[0]));
            StringBuilder code = new StringBuilder();
            String line;
            
            while ((line = reader.readLine()) != null) {
                code.append(line).append("\n");
            }
            reader.close();

            // Executa o código
            MaquinaPilha vm = new MaquinaPilha();
            vm.execute(code.toString().split("\n"));


        } catch (IOException e) {
            System.out.println("Erro ao ler arquivo: " + e.getMessage());
        } catch (Exception e) {
            System.out.println("Erro de execução: " + e.getMessage());
        }
    }



}