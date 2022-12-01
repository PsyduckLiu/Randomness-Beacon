import numpy as np
import matplotlib.pyplot as plt
import math

xArray = np.array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20])
data4 = np.loadtxt(open("time4.csv","rb"),delimiter=",")
data7 = np.loadtxt(open("time7.csv","rb"),delimiter=",")
data10 = np.loadtxt(open("time10.csv","rb"),delimiter=",")
data13 = np.loadtxt(open("time13.csv","rb"),delimiter=",")
data16 = np.loadtxt(open("time16.csv","rb"),delimiter=",")
data19 = np.loadtxt(open("time19.csv","rb"),delimiter=",")
data22 = np.loadtxt(open("time22.csv","rb"),delimiter=",")
data25 = np.loadtxt(open("time25.csv","rb"),delimiter=",")

data4 = data4[:-1]
data7 = data7[:-1]
data10 = data10[:-1]
data13 = data13[:-1]
data16 = data16[:-1]
data19 = data19[:-1]
data22 = data22[:-1]
data25 = data25[:-1]

plt.figure(1)
plt.xlabel('Round')
plt.ylabel('time(s)')
plt.xlim(1, 20)
plt.ylim(20, 50)
plt.plot(xArray, data4, "o")
plt.plot(xArray, data7, "o")
plt.plot(xArray, data10, "o")
plt.plot(xArray, data13, "o")
plt.plot(xArray, data16, "o")
plt.plot(xArray, data19, "o")
plt.plot(xArray, data22, "o")
plt.plot(xArray, data25, "o")
plt.plot(xArray, data4)
plt.plot(xArray, data7)
plt.plot(xArray, data10)
plt.plot(xArray, data13)
plt.plot(xArray, data16)
plt.plot(xArray, data19)
plt.plot(xArray, data22)
plt.plot(xArray, data25)
plt.title('title')
# plt.legend(('4','7','10','13','16','19','22','25'))
plt.savefig("consensus.png")

xArray = np.array([1, 2, 3, 4, 5, 6, 7, 8])
data4Mean = np.mean(data4)
data7Mean = np.mean(data7)
data10Mean = np.mean(data10)
data13Mean = np.mean(data13)
data16Mean = np.mean(data16)
data19Mean = np.mean(data19)
data22Mean = np.mean(data22)
data25Mean = np.mean(data25)
dataMean = np.array([data4Mean,data7Mean,data10Mean,data13Mean,data16Mean,data19Mean,data22Mean,data25Mean])

z1 = np.polyfit(xArray, dataMean, 1)
p1 = np.poly1d(z1)
print(z1)
print(p1)

plt.figure(2)
plt.xlabel('Number of Entropy Nodes')
plt.ylabel('Average Time(s)')
# plt.xlim(0, 70)
plt.plot(xArray, dataMean, "o")
plt.plot(xArray, dataMean)
plt.title('title')
plt.savefig("consensusAverage.png")