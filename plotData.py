import pandas as pd
import matplotlib.pyplot as plt

ax = plt.gca()

df = pd.read_csv('time_plot_canonical_vs_rapid.csv')

    
print(df)

df.plot(kind='scatter',marker="x",x='taxa',y='rapidnj',ax=ax, label="rapidnj")
df.plot(kind='scatter',marker="x",x='taxa',y='canonical', color='red', ax=ax, label="canonical")


plt.show()