import pandas as pd
import matplotlib.pyplot as plt

ax = plt.gca()

df = pd.read_csv('time_plot.csv')

df.plot(kind='line',x='taxa',y='rapidnj',ax=ax)
df.plot(kind='line',x='taxa',y='canonical', color='red', ax=ax)

plt.show()