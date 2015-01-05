# Create a "demo" namespace
ip netns add demo
 
# create a "veth" pair
ip link add veth0 type veth peer name veth1
 
# and move one to the namespace
ip link set veth1 netns demo
 
# configure the interfaces (up + IP)
ip netns exec demo ip link set lo up
ip netns exec demo ip link set veth1 up
ip netns exec demo ip addr add 169.254.1.2/30 dev veth1
ip link set veth0 up
ip addr add 169.254.1.1/30 dev veth0

# make sure ip forwarding is enabled
echo 1 > /proc/sys/net/ipv4/ip_forward
# enable Internet access for the namespace, assuming you ran the previous example
iptables -t nat -A POSTROUTING -i veth0 -j  MASQUERADE
# Forward main ":80" to guest ":80"
iptables -t nat -A PREROUTING -d <your main ip>/32 -p tcp --dport 80 -j  DNAT --to-destination  169.254.1.2:80

