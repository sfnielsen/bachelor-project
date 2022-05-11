import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

def plotData():
    ax = plt.gca()

    df = pd.read_csv('time_plot_canonical_vs_rapid.csv')

        
    print(df)

    df.plot(kind='scatter',marker="x",x='taxa',y='rapidnj',ax=ax, label="rapidnj")
    df.plot(kind='scatter',marker="x",x='taxa',y='canonical', color='red', ax=ax, label="canonical")


    plt.show()

def plotErrorbar():
    ax = plt.gca()

    df = pd.read_csv('time_plot_canonical_vs_rapid.csv')

    plt.errorbar(df.taxa, np.log(df.canonical), yerr=df.canonical_error, marker='x', label = 'CanonicalNJ',
                ecolor='red', fmt='None', capsize=2)
    plt.errorbar(df.taxa, np.log(df.rapidnj), yerr=df.rapidnj_error, marker='x', label = 'RapidNJ',
                ecolor='blue', fmt='None', capsize=2)
    plt.legend(loc ='upper left')

    plt.xlabel("# taxa")
    plt.ylabel("Waittime in MS (ln scale)")

    plt.show()

def plotAllTreesErrorbar():
    ax = plt.gca()
    
    df = pd.read_csv('allTrees_timetest.csv')
    plt.errorbar(df.taxa,  np.log(df.Norm), yerr= df.norm_err, marker='x', label = 'Norm',
                ecolor='blue', fmt='None', capsize=2)
    plt.errorbar(df.taxa,  np.log(df.Cluster_norm), yerr= df.cluster_err, marker='x', label = 'Cluster',
                ecolor='orange', fmt='None', capsize=2)
    plt.errorbar(df.taxa,  np.log(df.Spike_norm), yerr= df.spike_err, marker='x', label = 'Spike',
                ecolor='green', fmt='None', capsize=2)

    plt.legend(loc ='upper left')
    plt.xlabel("# taxa")
    plt.ylabel("Y axis label")

    plt.show()

def plotInitialRapidnjVsUUPDATErapidnj():
    ax = plt.gca()

    df_old = pd.read_csv('version_0_time.csv')
    df_new = pd.read_csv('version_4_time.csv')

    plt.errorbar(df_old.taxa, np.log(df_old.rapidnj), yerr=df_old.rapidnj_error, marker='x', label = 'RapidNJ_v0',
                ecolor='red', fmt='None', capsize=2)
    plt.errorbar(df_new.taxa, np.log(df_new.rapidnj), yerr=df_new.rapidnj_error, marker='x', label = 'RapidNJ_v4',
                ecolor='blue', fmt='None', capsize=2)
    plt.legend(loc ='upper left')

    plt.xlabel("# taxa")
    plt.ylabel("Waittime in MS (ln scale)")

    plt.show()

def plotScatter():
    array = pd.read_csv('cluster_bigmax.csv')
    xs = array.iloc[:, 0].to_list()
    ys = array.iloc[:,1:]
    for i in range(0,len(ys)):
        for j in range(0,len(ys.iloc[0])):
            plt.scatter(xs[i], np.log(ys.iloc[i,j]), marker='.', color='blue')

    plt.scatter(xs[len(ys)-1], np.log(ys.iloc[len(ys)-1,len(ys.iloc[0])-1]), marker='.', color='blue',label="Cluster")

    array1 = pd.read_csv('spike_bigmax.csv')
    xs1 = array1.iloc[:, 0].to_list()
    ys1 = array1.iloc[:,1:]
    for i in range(0,len(ys1)):
        for j in range(0,len(ys1.iloc[0])):
            plt.scatter(xs1[i], np.log(ys1.iloc[i,j]), marker='.', color='red')

    plt.scatter(xs1[i], np.log(ys1.iloc[i,j]), marker='.', color='red', label='Spike')


    df_new = pd.read_csv('version_5_time.csv')
    plt.errorbar(df_new.taxa, np.log(df_new.rapidnj), yerr=df_new.rapidnj_error, marker='x', label = 'RapidNJ_v5',
                ecolor='green', fmt='None', capsize=2)
    plt.legend(loc ='upper left')
    plt.xlabel("# taxa")
    plt.ylabel("Waittime in MS (ln scale)")

    plt.show()

plotScatter()