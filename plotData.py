import pandas as pd
import matplotlib.pyplot as plt

ax = plt.gca()

df = pd.read_csv('time_plot.csv')
print(df)
df.plot(kind='line',x='taxa',y='rapidnj_shifted',ax=ax)
df.plot(kind='line',x='taxa',y='canonical_shifted', color='red',ax=ax)
df.plot(kind='line',x='taxa',y='rapid_norm', color='green', ax=ax)

plt.show()