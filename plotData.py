import pandas as pd
import matplotlib.pyplot as plt

ax = plt.gca()

df = pd.read_csv('allTrees_timetest2.csv')

    
print(df)

df.plot(kind='line',marker="x",x='taxa',y='Sh_norm',ax=ax, label="Sh_norm")
df.plot(kind='line',marker="x",x='taxa',y='Norm', color='red', ax=ax, label="Norm")
df.plot(kind='line',marker="x",x='taxa',y='Uniform', color='purple', ax=ax, label="Uniform")
df.plot(kind='line',marker="x",x='taxa',y='Cluster_norm', color='green', ax=ax, label="Cluster_norm")
df.plot(kind='line',marker="x",x='taxa',y='Spike_norm', color='black', ax=ax, label="Spike_norm")


plt.show()