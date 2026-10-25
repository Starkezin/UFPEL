#include <iostream>
#include <fstream>
#include <string>
#include <vector>
#include <algorithm>
#include <iterator>
#include <chrono>
using namespace std;

vector<int> convertVecStrToVecInt(vector<string> vecStrings)
{
    vector<int> vecInt;
    for (auto str : vecStrings)
    {
        int aux = atoi(str.c_str());
        vecInt.push_back(aux);
    }
    return vecInt;
}
int calcTSP(vector<int> route, vector<vector<int>> distances)
{
    int total_distance = 0;
    for (int i = 0; i < route.size() - 1; i++)
    {
        total_distance += distances[route[i]][route[i + 1]];
    }
    total_distance += distances[route[route.size() - 1]][route[0]];
    return total_distance;
}
vector<string> getTokens(string line)
{
    int pos = 0;
    string token;
    vector<string> vecStrings;
    for (pos = 0; pos < line.size(); ++pos)
    {
        if (line[pos] != ' ')
        {
            token.insert(token.end(), line[pos]);
        }
        else
        {
            if (!token.empty())
            {
                vecStrings.push_back(token);
                token = "";
            }
        }
    }
    vecStrings.push_back(token);
    return vecStrings;
}

int main(int argc, char** argv)
{
    string filename = argv[1]; 

    ifstream fileInput(filename);
    string line;
    vector<vector<int>> distances;
    if (!fileInput)
    {
        cout << "Falha!" << endl;
    }               
    while (getline(fileInput, line))
    {
        vector<string> auxString = getTokens(line);
        vector<int> auxInt = convertVecStrToVecInt(auxString);
        distances.push_back(auxInt);
    }

    vector<int> list;
    int max_distance = 99999999;
    int curr_distance;
    for (int n = 1; n < distances.size(); n++)
    {
        list.push_back(n);
    }

    auto inicio = std::chrono::high_resolution_clock::now();
    do
    {
        vector<int> cr;
        cr.resize(list.size() + 1);
        cr[0] = 0;
        for (size_t i = 0; i < list.size(); ++i)
        {
            cr[i + 1] = list[i];
        }
        curr_distance = calcTSP(cr, distances);
        if (curr_distance < max_distance)
            max_distance = curr_distance;
    } while (next_permutation(list.begin(), list.begin() + distances.size() - 1));

    auto fim = std::chrono::high_resolution_clock::now();
    std::chrono::duration<double, std::milli> duracao = fim - inicio;
    cout << "Tempo de execucao: " << duracao.count() << " ms" << endl;
    cout << "final:" << max_distance << endl;
    return 0;
}
