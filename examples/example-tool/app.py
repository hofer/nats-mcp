# app.py
from mcp.server.fastmcp import FastMCP


mcp = FastMCP("OrderData")


@mcp.tool()
def query_order_data(question) -> str:
    """
    This tool can be used to query order data. The users question is used as a plain input
    into this tool. It can answer all queries about order ids, delivery/shipment locations, order
    value, total order value for a given timeframe, how much revenue in a certain timeframe etc.
    """
    return "some order data"


if __name__ == "__main__":
    #print(query_data("list 10 order ids"))
    mcp.run(transport="stdio")
