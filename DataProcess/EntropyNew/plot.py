import numpy as np
import matplotlib.pyplot as plt
import math

xArray = np.array([1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11])
data10 = np.loadtxt(open("10time.csv","rb"),delimiter=",")
data20 = np.loadtxt(open("20time.csv","rb"),delimiter=",")
data30 = np.loadtxt(open("30time.csv","rb"),delimiter=",")
data40 = np.loadtxt(open("40time.csv","rb"),delimiter=",")
data50 = np.loadtxt(open("50time.csv","rb"),delimiter=",")
data60 = np.loadtxt(open("60time.csv","rb"),delimiter=",")

# data4 = data4[:-1]
# data7 = data7[:-1]
# data10 = data10[:-1]
# data13 = data13[:-1]
# data16 = data16[:-1]
# data19 = data19[:-1]
# data22 = data22[:-1]
# data25 = data25[:-1]

plt.figure(1)
plt.xlabel('Round')
plt.ylabel('time(s)')
# plt.xlim(1, 20)
# plt.ylim(0, 50)
plt.plot(xArray, data10, "o")
plt.plot(xArray, data20, "o")
plt.plot(xArray, data30, "o")
plt.plot(xArray, data40, "o")
plt.plot(xArray, data50, "o")
plt.plot(xArray, data60, "o")
plt.plot(xArray, data10)
plt.plot(xArray, data20)
plt.plot(xArray, data30)
plt.plot(xArray, data40)
plt.plot(xArray, data50)
plt.plot(xArray, data60)
plt.title('title')
plt.legend(('10','20','30','40','50','60'))
plt.savefig("consensus.png")

xArray = np.array([10, 20, 30, 40, 50, 60])
data10Mean = np.mean(data10)
data20Mean = np.mean(data20)
data30Mean = np.mean(data30)
data40Mean = np.mean(data40)
data50Mean = np.mean(data50)
data60Mean = np.mean(data60)
dataMean = np.array([data10Mean,data20Mean,data30Mean,data40Mean,data50Mean,data60Mean])

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
plt.savefig("entropyAverage.png")