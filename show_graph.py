import json
import networkx as nx
import matplotlib.pyplot as plt

def nudge(pos, x_shift, y_shift):
    return {n:(x + x_shift, y + y_shift) for n,(x,y) in pos.items()}

with open('graph.cyjs') as f:
    data = json.load(f)

G = nx.cytoscape_graph(data)

pos = nx.spring_layout(G, k=0.5, iterations=50)
pos_nodes = nudge(pos, 0, 0.05)
nx.draw_networkx(G, pos=pos, with_labels=False)
nx.draw_networkx_labels(G, pos=pos_nodes)
plt.show()

#nx.draw_networkx(G)
#
#ax = plt.gca()
#ax.margins(0.20)
#plt.axis("off")
#plt.show()
