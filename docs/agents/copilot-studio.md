# zara-jira-mcp — Microsoft Copilot Studio Setup

Copilot Studio supports MCP servers as agent extensions.

## Configuration

In Copilot Studio:

1. Open your agent > **Actions** > **Add an action**
2. Select **Model Context Protocol (MCP)**
3. Configure the connection:
   - **Transport:** stdio
   - **Command:** `zara-jira-mcp-wrapper`

## For Self-Hosted / On-Prem

If running the MCP server as a remote HTTP endpoint:

1. Deploy zara-jira-mcp with Streamable HTTP transport
2. In Copilot Studio, add the MCP server URL
3. Configure authentication as needed

## Notes

- Copilot Studio automatically discovers tools and creates actions from MCP server capabilities
- The 124 PM tools become available as actions in your Copilot agent
- Best for enterprise teams already on Microsoft 365
- See `docs/agents/README.md` for the wrapper script and environment setup
