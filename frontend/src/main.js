import "./style.css";
import "./app.css";
import {
  OpenNewInstance,
  GetStartupService,
  GetVersion,
  SaveWindowPositionManual,
} from "../wailsjs/go/main/App";
import { WindowSetTitle } from "../wailsjs/runtime/runtime";

const aiServices = [
  {
    id: "chatgpt",
    label: "ChatGPT",
    url: "https://chatgpt.com",
    description:
      "<b>Most popular general-purpose AI</b><br><br>" +
      "Powered by OpenAI's GPT-4 and GPT-5 models. Excels at creative writing, code generation, problem-solving, and conversational tasks. Fast response times with multimodal capabilities (text, images, voice).<br><br>" +
      "<b>Best for:</b> Content creation, coding assistance, learning, brainstorming, and everyday tasks.",
  },
  {
    id: "claude",
    label: "Claude (Sonnet)",
    url: "https://claude.ai",
    description:
      "<b>Deep reasoning and analysis</b><br><br>" +
      "Anthropic's Claude Sonnet excels at nuanced understanding, long-context analysis (200K+ tokens), and following complex instructions. Strong ethical guidelines and safety focus. Better at structured analysis than creative tasks.<br><br>" +
      "<b>Best for:</b> Document analysis, research synthesis, technical writing, code review, and ethical reasoning.",
  },
  {
    id: "copilot",
    label: "Copilot",
    url: "https://copilot.microsoft.com",
    description:
      "<b>Microsoft ecosystem integration</b><br><br>" +
      "Integrated with Microsoft 365 apps (Word, Excel, PowerPoint, Outlook). Combines GPT-4 with Bing search for grounded, up-to-date answers. Supports plugins and organizational data access with enterprise security.<br><br>" +
      "<b>Best for:</b> Office productivity, business workflows, enterprise tasks, and real-time web research.",
  },
  {
    id: "deepseek",
    label: "Deepseek",
    url: "https://chat.deepseek.com/",
    description:
      "<b>Advanced reasoning and coding</b><br><br>" +
      "Chinese open-source model (DeepSeek-V3.2) with strong mathematical and coding capabilities. Features chain-of-thought reasoning and competitive performance at lower costs. Newly enhanced with agent capabilities and thinking modes.<br><br>" +
      "<b>Best for:</b> Complex coding tasks, mathematical problem-solving, algorithmic challenges, and cost-effective AI access.",
  },
  {
    id: "gemini",
    label: "Gemini",
    url: "https://gemini.google.com",
    description:
      "<b>Google's multimodal powerhouse</b><br><br>" +
      "Latest Gemini 2.0 Flash and 2.5 Pro models with advanced multimodal understanding (text, images, video, audio). Deep integration with Google Workspace and Search. Excels at visual tasks, data analysis, and creative content.<br><br>" +
      "<b>Best for:</b> Image generation, video analysis, Google Workspace tasks, research with web grounding, and visual creativity.",
  },
  {
    id: "grok",
    label: "Grok",
    url: "https://grok.com",
    description:
      "<b>Real-time X/Twitter integration</b><br><br>" +
      "X's AI with direct access to real-time X/Twitter data and trending topics. More conversational and less filtered than competitors. Developed by xAI with focus on truthfulness and current events awareness.<br><br>" +
      "<b>Best for:</b> Social media insights, trending topics, current events, real-time news analysis, and uncensored conversations.",
  },
  {
    id: "meta",
    label: "Meta AI",
    url: "https://www.meta.ai",
    description:
      "<b>Social-first AI assistant</b><br><br>" +
      "Meta's LLaMA-powered AI integrated across Facebook, Instagram, and WhatsApp. Focuses on conversational AI, image generation, and social interactions. Privacy-conscious with transparent data usage policies.<br><br>" +
      "<b>Best for:</b> Social media content, casual conversations, image creation, and Facebook/Instagram-related tasks.",
  },
  {
    id: "perplexity",
    label: "Perplexity",
    url: "https://www.perplexity.ai",
    description:
      "<b>AI-powered research engine</b><br><br>" +
      "Combines conversational AI with real-time web search and citations. Every answer includes source links for verification. Excels at research, fact-checking, and providing up-to-date information with transparency.<br><br>" +
      "<b>Best for:</b> Academic research, fact-checking, current events, cited answers, and information discovery with sources.",
  },
];

let currentService = "chatgpt";

// Check if we should navigate to a specific service on startup
GetStartupService().then((startupService) => {
  if (startupService && startupService !== "") {
    // We were launched with a service argument, navigate to it
    const service = aiServices.find((s) => s.id === startupService);
    if (service) {
      WindowSetTitle(`SimpleAI - ${service.label}`);
      window.location.href = service.url;
      return;
    }
  }

  // No startup service, show launcher
  showLauncher();
  WindowSetTitle("SimpleAI");
});

function showLauncher() {
  // Show launcher with service buttons
  GetVersion().then((version) => {
    // Extract major.minor from version (e.g., "1.0.0" -> "1.0")
    const shortVersion = version.split(".").slice(0, 2).join(".");
    document.querySelector("#app").innerHTML = `
    <div style="
      --wails-draggable: drag;
      display: flex;
      flex-direction: column;
      height: 100vh;
      background: rgba(27, 38, 54, 1);
      overflow: hidden;
    ">
      <div id="titlebar" style="
        display: flex;
        align-items: center;
        justify-content: flex-end;
        height: 30px;
        background: rgba(27, 38, 54, 1);
        flex-shrink: 0;
      ">
        <span style="
          color: white;
          font-size: 13px;
          user-select: none;
          position: absolute;
          left: 10px;
        ">SimpleAI ${shortVersion}</span>
        <div style="
          display: flex;
          gap: 0px;
        ">
          <button id="btn-minimize" style="
            --wails-draggable: no-drag;
            background: none;
            border: none;
            color: white;
            font-size: 16px;
            width: 30px;
            height: 30px;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: all 0.3s ease;
            border-radius: 4px;
          " 
          onmouseover="this.style.background='rgba(0, 212, 255, 0.3)'; this.style.boxShadow='0 0 10px rgba(0, 212, 255, 0.5)'; this.style.transform='scale(1.1)';" 
          onmouseout="this.style.background='none'; this.style.boxShadow='none'; this.style.transform='scale(1)';"
          title="Minimize">-</button>
          <button id="btn-close" style="
            --wails-draggable: no-drag;
            background: none;
            border: none;
            color: white;
            font-size: 16px;
            width: 30px;
            height: 30px;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: all 0.3s ease;
            border-radius: 4px;
          " 
          onmouseover="this.style.background='rgba(255, 50, 80, 0.3)'; this.style.boxShadow='0 0 10px rgba(255, 50, 80, 0.5)'; this.style.transform='scale(1.1)';" 
          onmouseout="this.style.background='none'; this.style.boxShadow='none'; this.style.transform='scale(1)';"
          title="Close">×</button>
        </div>
      </div>
      <div style="
        --wails-draggable: no-drag;
        display: flex;
        flex-wrap: wrap;
        gap: 5px;
        justify-content: center;
        align-items: center;
      ">
        ${aiServices
          .map(
            (service) => `
          <div style="
            position: relative;
            width: 150px;
          ">
            <button id="btn-${service.id}" style="
              width: 100%;
              padding: 5px 5px;
              font-size: 16px;
              background: rgba(0, 212, 255, 0.1);
              border: 2px solid #00d4ff;
              color: white;
              border-radius: 10px;
              cursor: pointer;
              transition: all 0.2s;
              white-space: nowrap;
            " onmouseover="this.style.background='rgba(0, 212, 255, 0.2)'" 
               onmouseout="this.style.background='rgba(0, 212, 255, 0.1)'">
              ${service.label}
            </button>
            <button id="info-${service.id}" style="
              --wails-draggable: no-drag;
              position: absolute;
              top: 2px;
              right: 2px;
              width: 20px;
              height: 20px;
              border-radius: 50%;
              background: rgba(0, 212, 255, 0.3);
              border: 1px solid #00d4ff;
              color: white;
              font-size: 12px;
              cursor: pointer;
              display: flex;
              align-items: center;
              justify-content: center;
              transition: all 0.2s;
            " onmouseover="this.style.background='rgba(0, 212, 255, 0.5)'; this.style.transform='scale(1.1)';" 
               onmouseout="this.style.background='rgba(0, 212, 255, 0.3)'; this.style.transform='scale(1)';">
              ?
            </button>
          </div>
        `,
          )
          .join("")}
      </div>
      </div>
    </div>
  `;

    // Add click handlers to open new instances
    aiServices.forEach((service) => {
      document
        .getElementById(`btn-${service.id}`)
        .addEventListener("click", async () => {
          try {
            await OpenNewInstance(service.id);
          } catch (err) {
            console.error("Failed to open new instance:", err);
          }
        });

      // Add info icon click handler
      document
        .getElementById(`info-${service.id}`)
        .addEventListener("click", (e) => {
          e.stopPropagation();
          const desc = service.description || "No description available.";

          // Create custom tooltip/modal
          const modal = document.createElement("div");
          modal.style.cssText = `
            position: fixed;
            top: 30px;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(27, 38, 54, 0.98);
            z-index: 10000;
            color: white;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            padding: 10px;
            box-sizing: border-box;
          `;
          modal.innerHTML = `
            <div style="font-size: 16px; font-weight: bold; margin-bottom: 10px;">${service.label}</div>
            <div style="font-size: 12px; margin-bottom: 10px; text-align: left; max-width: 600px;">${desc}</div>
            <button id="close-modal" style="
              padding: 5px ;
              background: rgba(0, 212, 255, 0.2);
              border: 2px solid #00d4ff;
              color: white;
              border-radius: 10px;
              cursor: pointer;
              font-size: 16px;
              transition: all 0.2s;
            " onmouseover="this.style.background='rgba(0, 212, 255, 0.3)'" 
               onmouseout="this.style.background='rgba(0, 212, 255, 0.2)'">Close</button>
          `;

          document.body.appendChild(modal);

          document
            .getElementById("close-modal")
            .addEventListener("click", () => {
              document.body.removeChild(modal);
            });

          // Close on background click
          modal.addEventListener("click", (e) => {
            if (e.target === modal) {
              document.body.removeChild(modal);
            }
          });
        });
    });

    // Add window control handlers
    document.getElementById("btn-minimize").addEventListener("click", () => {
      window.runtime.WindowMinimise();
    });

    document.getElementById("btn-close").addEventListener("click", () => {
      window.runtime.Quit();
    });
  });
}

// Save window position on resize/move events with debouncing
let savePositionTimeout = null;
const savePositionDebounced = () => {
  if (savePositionTimeout) {
    clearTimeout(savePositionTimeout);
  }
  savePositionTimeout = setTimeout(async () => {
    try {
      await SaveWindowPositionManual();
    } catch (err) {
      console.error("Failed to save window position:", err);
    }
  }, 500); // Wait 500ms after last resize/move
};

window.addEventListener("resize", savePositionDebounced);
// Note: 'move' event doesn't exist in standard DOM, position changes handled on close
