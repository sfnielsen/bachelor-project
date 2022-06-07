import pandas as pd
import matplotlib.pyplot as plt
import numpy as np
import seaborn as sns


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

def plotHeatMap(name):
    dataframe = pd.read_csv(name)
    nparray = dataframe.to_numpy()
    plt.xlabel("column # in S")
    plt.ylabel("row # in S")
    plt.imshow(nparray, cmap='hot', interpolation='nearest')
    plt.show()

def findTotalLookups():
    spike = pd.read_csv('s_spike_lookup_analysis_400.csv')
    nparray = spike.to_numpy()
    sum = nparray.sum()
    print("spike_lookup sum: ", sum)

    clust = pd.read_csv('s_clus_lookup_analysis_400.csv')
    nparray_clust = clust.to_numpy()
    sum1 = nparray_clust.sum()
    print("clust_lookup sum: ", sum1)




    spike = pd.read_csv('s_spike_update_analysis_400.csv')
    nparray = spike.to_numpy()
    sum = nparray.sum()
    print("spike_update sum: ", sum)

    clust = pd.read_csv('s_clus_update_analysis_400.csv')
    nparray_clust = clust.to_numpy()
    sum1 = nparray_clust.sum()
    print("clust_update sum: ", sum1)


def plotDifferentTimes():
    ax = plt.gca()

    normal = pd.read_csv('normal_localtimes.csv')
    cluster = pd.read_csv('cluster_localtimes_onego.csv')
    cluster = cluster - [0, 11082483456,31942606926,406299230494]
    spike = pd.read_csv('spike_localtimes_onego.csv')


    normal['normal_updates'] = np.log(normal['normal_updates'])
    normal['normal_inits'] =   np.log(normal['normal_inits'])
    normal['normal_lookups'] = np.log(normal['normal_lookups'])

    normal.plot(kind='scatter',marker="x",x='taxa',y='normal_inits',ax=ax,color="red", label="initialization normal")
    normal.plot(kind='scatter',marker="x",x='taxa',y='normal_updates',ax=ax, color="blue", label="update sorting normal")
    normal.plot(kind='scatter',marker="x",x='taxa',y='normal_lookups',ax=ax, color="green", label="# lookups normal")


    cluster['normal_updates'] = np.log(cluster['normal_updates'])
    cluster['normal_inits'] =   np.log(cluster['normal_inits'])
    cluster['normal_lookups'] = np.log(cluster['normal_lookups'])

    cluster.plot(kind='scatter',marker=".",x='taxa',y='normal_inits',ax=ax,color="red", label="initialization cluster")
    cluster.plot(kind='scatter',marker=".",x='taxa',y='normal_updates',ax=ax, color="blue", label="update sorting cluster")
    cluster.plot(kind='scatter',marker=".",x='taxa',y='normal_lookups',ax=ax, color="green", label="# lookups cluster")

    spike['normal_updates'] = np.log(spike['normal_updates'])
    spike['normal_inits'] =   np.log(spike['normal_inits'])
    spike['normal_lookups'] = np.log(spike['normal_lookups'])

    spike.plot(kind='scatter',marker="o",x='taxa',y='normal_inits',ax=ax,color="red", label="initialization spike")
    spike.plot(kind='scatter',marker="o",x='taxa',y='normal_updates',ax=ax, color="blue", label="update sorting spike")
    spike.plot(kind='scatter',marker="o",x='taxa',y='normal_lookups',ax=ax, color="green", label="# lookups spike")
    plt.show()


def plotDifferentTimes():
    ax = plt.gca()

    cluster = pd.read_csv('cluster_localtimes_onego.csv')
    taxa = cluster['taxa']
    cluster = cluster - [0, 11082483456,31942606926,406299230494]
    cluster = cluster.diff()
    cluster['taxa'] = taxa
    spike = pd.read_csv('spike_localtimes_onego.csv')
    spike = spike.diff()
    spike['taxa'] = taxa

    cluster['normal_updates'] = np.log(cluster['normal_updates'])
    cluster['normal_inits'] =   np.log(cluster['normal_inits'])
    cluster['normal_lookups'] = np.log(cluster['normal_lookups'])

    cluster.plot(kind='scatter',marker=".",x='taxa',y='normal_inits',ax=ax,color="red", label="initialization cluster")
    cluster.plot(kind='scatter',marker=".",x='taxa',y='normal_updates',ax=ax, color="blue", label="update sorting cluster")
    cluster.plot(kind='scatter',marker=".",x='taxa',y='normal_lookups',ax=ax, color="green", label="# lookups cluster")

    spike['normal_updates'] = np.log(spike['normal_updates'])
    spike['normal_inits'] =   np.log(spike['normal_inits'])
    spike['normal_lookups'] = np.log(spike['normal_lookups'])

    spike.plot(kind='scatter',marker="o",x='taxa',y='normal_inits',ax=ax,color="red", label="initialization spike")
    spike.plot(kind='scatter',marker="o",x='taxa',y='normal_updates',ax=ax, color="blue", label="update sorting spike")
    spike.plot(kind='scatter',marker="o",x='taxa',y='normal_lookups',ax=ax, color="green", label="# lookups spike")


    cluster2 = pd.read_csv('cluster_times_lastrow_heuristic.csv')
    spike2 = pd.read_csv('spike_times_lastrow_heuristic.csv')

    cluster2['normal_updates'] = np.log(cluster2['normal_updates'])
    cluster2['normal_inits'] =   np.log(cluster2['normal_inits'])
    cluster2['normal_lookups'] = np.log(cluster2['normal_lookups'])

    cluster2.plot(kind='scatter',marker="x",x='taxa',y='normal_inits',ax=ax,color="red", label="initialization cluster2")
    cluster2.plot(kind='scatter',marker="x",x='taxa',y='normal_updates',ax=ax, color="blue", label="update sorting cluster2")
    cluster2.plot(kind='scatter',marker="x",x='taxa',y='normal_lookups',ax=ax, color="green", label="# lookups cluster2")

    spike2['normal_updates'] = np.log(spike2['normal_updates'])
    spike2['normal_inits'] =   np.log(spike2['normal_inits'])
    spike2['normal_lookups'] = np.log(spike2['normal_lookups'])

    spike2.plot(kind='scatter',marker="^",x='taxa',y='normal_inits',ax=ax,color="red", label="initialization spike2")
    spike2.plot(kind='scatter',marker="^",x='taxa',y='normal_updates',ax=ax, color="blue", label="update sorting spike2")
    spike2.plot(kind='scatter',marker="^",x='taxa',y='normal_lookups',ax=ax, color="green", label="# lookups spike2")

    plt.show()




plotDifferentTimes()