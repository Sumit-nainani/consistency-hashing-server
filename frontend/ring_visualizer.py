import matplotlib.pyplot as plt
import numpy as np

def draw_ring(node_positions=None, request_position=None):
    """
    Draws a consistent hashing ring with 64 positions.
    Optional:
    - node_positions: list of ints (positions of servers)
    - request_position: int (position of request)
    """
    ring_size = 64
    radius = 1.0
    angles = np.linspace(0, 2 * np.pi, ring_size, endpoint=False)

    fig, ax = plt.subplots(figsize=(6, 6), subplot_kw={'polar': True})
    ax.set_xticks(angles)
    ax.set_yticklabels([])
    ax.set_xticklabels([str(i) for i in range(ring_size)], fontsize=7)
    ax.set_theta_zero_location("N")
    ax.set_theta_direction(-1)

    # Draw only the border ring
    ax.plot(angles, [radius] * ring_size, 'o', color='gray', markersize=4)

    # Highlight node positions (servers)
    if node_positions:
        for pos in node_positions:
            angle = angles[pos % ring_size]
            ax.plot(angle, radius, 'o', color='blue', markersize=10)

    # Highlight request position
    if request_position is not None:
        angle = angles[request_position % ring_size]
        ax.plot(angle, radius, 'o', color='red', markersize=10)

    ax.set_title("Consistent Hashing Ring", va='bottom')
    return fig