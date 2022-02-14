from Bio import Phylo
from nbformat import read
from io import StringIO

myfile = open("newick.txt", "r")
newick=myfile.read()
    
input = StringIO(newick)
tree = Phylo.read(input, "newick")
print(tree)
Phylo.draw(tree)  