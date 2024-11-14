import path from "path";
import ejs from "ejs";

export const handler = async (event: any) => {
  try {
    const templatePath = path.join(__dirname, "../views/template.ejs");
    const html = await ejs.renderFile(templatePath, {
      title: "My Page",
      content: "Hello, world!",
    });

    return {
      statusCode: 200,
      body: html,
      headers: { "Content-Type": "text/html" },
    };
  } catch (error) {
    console.error(error);
    return {
      statusCode: 500,
      body: "Internal Server Error",
    };
  }
};
