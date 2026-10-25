package ast;

import java.util.ArrayList;

public class Prog{
    public Main main;
    public ArrayList<Fun> fun;
    public Prog(Main main, ArrayList<Fun> fun)
    {
        this.main = main;
        this.fun = fun;
    }

    public String geraCodigo() {
        StringBuilder sb = new StringBuilder();
        
        // 1. Cabeçalhos (Apenas uma vez aqui no topo)
        sb.append("#include <stdio.h>\n");
        sb.append("#include <stdlib.h>\n\n");
        
        // 2. Funções (Devem vir ANTES da main em C)
        if (fun != null && !fun.isEmpty()) {
            for (Fun f : fun) {
                sb.append(f.geraCodigo());
            }
        }
        
        // 3. Main
        if (main != null) {
            sb.append(main.geraCodigo());
        }
        
        return sb.toString();
    }
}

