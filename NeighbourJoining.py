import numpy as np   
from io import StringIO
from Bio import Phylo

def main(D,labels):
    
    if len(D) == len(labels):
        neighbour_join(D,labels)

def neighbour_join(D,labels):
    n = len(D)
    M = np.zeros((n,n))
    r = np.zeros(n)
    r = np.sum(D, axis= 1) / (n-2) #rowwise summation
    for i in range(0,n):
        for j in range(i,n):
            if i == j:
                M[i,j] = 0
            else:
                M[i,j] = D[i,j] - (r[i]+r[j]) 
    
    pair_ij = np.where(M== np.min(M)) 
    #if multiple we just pick one pair
    p_i = pair_ij[0][0]
    p_j =pair_ij[1][0]
    
    #i and j distances to new node u
    v_iu = D[p_i,p_j]/2 + (r[p_i]-r[p_j])/2
    v_ju = D[p_i,p_j]/2 + (r[p_j]-r[p_i])/2
    #creating newick form
    labels[p_i] = ("(" + labels[p_i] + ":" + str(v_iu) + "," + labels[p_j] + ":" + str(v_ju) + ")")
    labels = np.delete(labels, p_j)
    
    D_new, new_labels = compute_new_dist_mat(D, p_i,p_j, labels)

    if len(D_new) > 2:
        neighbour_join(D_new,new_labels)
    else:
        #creating newick form
        newick = "(" + new_labels[p_i] + ":" + str(D_new[p_i,p_j] / 2 ) +  "," + new_labels[p_j] + ":" + str(D_new[p_i,p_j] / 2 )+ ");"
        print(newick)
        input = StringIO(newick)
        tree = Phylo.read(input, "newick")
        print(tree)
        Phylo.draw(tree)
        


def compute_new_dist_mat(D,p_i ,p_j, labels):
    n = len(D)
    for k in range(0, n):
        if p_i == k:
            D[p_i,k] = 0
        if p_j == k:
            continue
        else:
            #overwrite p_i as merge ij
            val = (D[p_i,k] + D[p_j,k] - D[p_i,p_j]) / 2 
            D[p_i,k] = val
            D[k,p_i] = val
    #delete p_j
    D = np.delete(D,[p_j], axis = 0)
    D = np.delete(D,[p_j], axis = 1)
    
    return D, labels

if __name__ == "__main__":
    
    D = np.array(
        [
            [  0,  5,  68,  57, 127,  27,  28,  33],
            [ 5,   0,  58,  47, 117,  8,  52,  57],
            [ 68,  58,   0,  35,  69,  35,  87,  92],
            [ 57,  47,  35,   0,  94,  44,  79,  84],
            [127, 117,  69,  94,   0, 144, 149, 154],
            [ 27,  8,  35,  44, 144,   0,  27,  54],
            [ 28,  52,  87,  79, 149,  27,   0,  13],
            [ 33,  57,  92,  84, 154,  54,  13,   0],
         ]
        )
    
    labels = np.array(["A","B","C","D", "E", "F", "G", "H"], type(str))
    """
    D = np.array([[0,17,21,27],
                  [17,0,12,18],
                  [21,12,0,14],
                  [27,18,14,0]])
    labels = np.array(["A","B","C","D"], type(str))
    """
        
    main(D,labels)