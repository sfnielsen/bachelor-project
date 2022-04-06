import pandas as pd
import matplotlib.pyplot as plt

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

    plt.errorbar(df.taxa, df.canonical, yerr=df.canonical_error, marker='x', label = 'CanonicalNJ',
                ecolor='red', fmt='None', capsize=2)
    plt.errorbar(df.taxa, df.rapidnj, yerr=df.rapidnj_error, marker='x', label = 'RapidNJ',
                ecolor='blue', fmt='None', capsize=2)
    plt.legend(loc ='upper left')

    plt.show()

plotErrorbar()