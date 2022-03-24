import json
import networkx as nx
import matplotlib.pyplot as plt

def nudge(pos, x_shift, y_shift):
    return {n:(x + x_shift, y + y_shift) for n,(x,y) in pos.items()}

with open('graph.cyjs') as f:
    data = json.load(f)

G = nx.cytoscape_graph(data)

# Remove the alarm nodes, to improve visualization.
# Create a copy of the graph, to avoid modifying while iterating
for node in G.copy().nodes():
    # Drop alarm nodes, its ids are only numeric
    try:
        int(node)
        G.remove_node(node)
    except ValueError:
        pass

pos = nx.spring_layout(G, k=0.5, iterations=50)
pos_nodes = nudge(pos, 0, 0.05)
nx.draw_networkx(G, pos=pos, with_labels=False)
nx.draw_networkx_labels(G, pos=pos_nodes)
plt.show()
