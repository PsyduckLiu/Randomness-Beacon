import numpy as np
import matplotlib.pyplot as plt
import math

data0 = np.loadtxt(open("resultCloud2.csv","rb"),delimiter=",",usecols=[0,1])
data1 = np.loadtxt(open("resultWSL3.csv","rb"),delimiter=",",usecols=[0,1])
data2 = np.loadtxt(open("resultWSLASUS1.csv","rb"),delimiter=",",usecols=[0,1])

x0 = data0[:,1]
x0 = x0.astype('int')
print('x is :\n',x0)

time0 = data0[:,0]
time0 = time0.astype('float')
print('time is :\n',time0)

log_time0 = np.log(time0)
coefficients0 = np.polyfit(x0, log_time0, 1)
print(coefficients0)

timeVals0 = np.exp(coefficients0[1]) * np.exp(coefficients0[0]*x0)
print(math.exp(coefficients0[1]),"*",math.exp(coefficients0[0]),"^x")

x1 = data1[:,1]
x1 = x1.astype('int')
print('x is :\n',x1)

time1 = data1[:,0]
time1 = time1.astype('float')
print('time is :\n',time1)

log_time1 = np.log(time1)
coefficients1 = np.polyfit(x1, log_time1, 1)
print(coefficients1)

timeVals1 = np.exp(coefficients1[1]) * np.exp(coefficients1[0]*x1)
print(math.exp(coefficients1[1]),"*",math.exp(coefficients1[0]),"^x")

x2 = data2[:,1]
x2 = x2.astype('int')
print('x is :\n',x2)

time2 = data2[:,0]
time2 = time2.astype('float')
print('time is :\n',time2)

log_time2 = np.log(time2)
coefficients2 = np.polyfit(x2, log_time2, 1)
print(coefficients2)

timeVals2 = np.exp(coefficients2[1]) * np.exp(coefficients2[0]*x2)
print(math.exp(coefficients2[1]),"*",math.exp(coefficients2[0]),"^x")

plt.xlabel('k')
plt.ylabel('time(s)')
plt.plot(x0, time0, "o")
plt.plot(x0, timeVals0)
plt.plot(x1, time1, "o")
plt.plot(x1, timeVals1)
plt.plot(x2, time2, "o")
plt.plot(x2, timeVals2)
plt.title('title')
plt.legend(('origin','fit'))
plt.savefig("origin.png")
