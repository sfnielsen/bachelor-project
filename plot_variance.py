# Implementation of matplotlib function
import numpy as np
import matplotlib.pyplot as plt
   
fig = plt.figure()
x = np.arange(10)
y = 3 * np.sin(x / 20 * np.pi)
yerr = np.linspace(0.05, 0.2, 10)
print(yerr)
   
plt.errorbar(x, y+7, yerr = yerr,
             label ='Line1', ecolor='red',fmt='None', capsize=2)


plt.errorbar(x, y + 5, yerr = yerr,
             uplims = True, 
             label ='Line2')
plt.errorbar(x, y + 3, yerr = yerr, 
             uplims = True, 
             lolims = True,
             label ='Line3')
  
upperlimits = [True, False] * 5
lowerlimits = [False, True] * 5
plt.errorbar(x, y, yerr = yerr,
             uplims = upperlimits, 
             lolims = lowerlimits,
             label ='Line4')
   
plt.legend(loc ='upper left')
  
plt.title('matplotlib.pyplot.errorbar()\
function Example')
plt.show()