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

    df_old = pd.read_csv('time_plot_canonical_vs_rapid.csv')[0:24]
    df_new = pd.read_csv('time_plot_canonical_vs_rapid_from_u_update.csv')

    plt.errorbar(df_old.taxa, np.log(df_old.rapidnj), yerr=df_old.rapidnj_error, marker='x', label = 'RapidNJ_intial_version',
                ecolor='red', fmt='None', capsize=2)
    plt.errorbar(df_old.taxa, np.log(df_new.rapidnj), yerr=df_new.rapidnj_error, marker='x', label = 'RapidNJ_New',
                ecolor='blue', fmt='None', capsize=2)
    plt.legend(loc ='upper left')

    plt.xlabel("# taxa")
    plt.ylabel("Waittime in MS (ln scale)")

    plt.show()

plotInitialRapidnjVsUUPDATErapidnj()