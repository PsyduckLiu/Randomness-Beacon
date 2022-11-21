import numpy as np
import matplotlib.pyplot as plt
import math

# data = np.loadtxt(open("resultCloud1.csv","rb"),delimiter=",",usecols=[0,1])
data = np.loadtxt(open("resultWSL1.csv","rb"),delimiter=",",usecols=[0,1])

x = data[:,1]
x = x.astype('int')
print('x is :\n',x)

time = data[:,0]
time = time.astype('float')
print('time is :\n',time)

log_time = np.log(time)
coefficients = np.polyfit(x, log_time, 1)
print(coefficients)

timeVals = np.exp(coefficients[1]) * np.exp(coefficients[0]*x)
print(math.exp(coefficients[1]),"*",math.exp(coefficients[0]),"^x")

plt.xlabel('k')
plt.ylabel('time(s)')
plt.plot(x, time, "o")
plt.plot(x, timeVals)
plt.title('title')
plt.legend(('origin','fit'))
plt.savefig("origin.png")
