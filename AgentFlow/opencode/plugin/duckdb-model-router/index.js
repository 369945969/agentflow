export default async function DuckDBModelRouterPlugin(ctx) {
  return {
    "chat.message": async () => {},
    meta: {
      directory: ctx?.directory,
    },
  };
}
