import{u as $,r as f,j as e,P as m,c as x,aj as O,aW as A}from"./index-C6MQE079.js";import{u as k,B as T}from"./index-C5I3Iaw-.js";import{c as N}from"./index-DZDPmYqs.js";import{u as P}from"./theme-context-wCTkV0aR.js";import{M as U}from"./index-B0HEP2-Y.js";import{T as R}from"./index-81ZfuQxj.js";import{u as D}from"./app-context-BaKmc3Lv.js";const S="_option_1gdw6_1",M="_active_1gdw6_6",V="_iframeIcon_1gdw6_9",F="_scriptsIcon_1gdw6_12",Y="_chromePluginIcon_1gdw6_15",B="_pluginInstallIcon_1gdw6_18",p={option:S,active:M,iframeIcon:V,scriptsIcon:F,chromePluginIcon:Y,pluginInstallIcon:B},d={iframe:{getContent:(t,r)=>`<iframe
 src="${t}${m}/chatbot/${r}"
 style="width: 100%; height: 100%; min-height: 700px"
 frameborder="0"
 allow="microphone">
</iframe>`},scripts:{getContent:(t,r,h,c)=>`<script>
 window.gofyChatbotConfig = {
  token: '${r}'${c?`,
  isDev: true`:""}${A?`,
  baseUrl: '${t}${m}'`:""},
  inputs: {
    // You can define the inputs from the Start node here
    // key is the variable name
    // e.g.
    // name: "NAME"
  },
  systemVariables: {
    // user_id: 'YOU CAN DEFINE USER ID HERE',
    // conversation_id: 'YOU CAN DEFINE CONVERSATION ID HERE, IT MUST BE A VALID UUID',
  },
  userVariables: {
    // avatar_url: 'YOU CAN DEFINE USER AVATAR URL HERE',
    // name: 'YOU CAN DEFINE USER NAME HERE',
  },
 }
<\/script>
<script
 src="${t}${m}/embed.min.js"
 id="${r}"
 defer>
<\/script>
<style>
  #gofy-chatbot-bubble-button {
    background-color: ${h} !important;
  }
  #gofy-chatbot-bubble-window {
    width: 24rem !important;
    height: 40rem !important;
  }
</style>`},chromePlugin:{getContent:(t,r)=>`ChatBot URL: ${t}${m}/chatbot/${r}`}},a="overview.appInfo.embedded",K=({siteInfo:t,isShow:r,onClose:h,appBaseUrl:c,accessToken:u,className:j})=>{var y;const{t:n}=$(),[s,C]=f.useState("iframe"),[l,b]=f.useState({iframe:!1,scripts:!1,chromePlugin:!1}),{odysseythinkVersionInfo:v}=D(),g=P();g.buildTheme((t==null?void 0:t.chat_color_theme)??null,(t==null?void 0:t.chat_color_theme_inverted)??!1);const w=v.current_env==="TESTING"||v.current_env==="DEVELOPMENT",_=()=>{var o;if(s==="chromePlugin"){const i=d[s].getContent(c,u).split(": ");i.length>1&&N(i[1])}else N(d[s].getContent(c,u,((o=g.theme)==null?void 0:o.primaryColor)??"#1C64F2",w));b({...l,[s]:!0})},E=()=>{const o={...l};Object.keys(o).forEach(i=>{o[i]=!1}),b(o)},I=()=>{window.open("https://chrome.google.com/webstore/detail/gofy-chatbot/ceehdapohffmjmkdcifjofadiaoeggaf","_blank","noopener,noreferrer")};return f.useEffect(()=>{E()},[r]),e.jsxs(U,{title:n(`${a}.title`,{ns:"appOverview"}),isShow:r,onClose:h,className:"w-[640px] !max-w-2xl",wrapperClassName:j,closable:!0,children:[e.jsx("div",{className:"system-sm-medium mb-4 mt-8 text-text-primary",children:n(`${a}.explanation`,{ns:"appOverview"})}),e.jsx("div",{className:"flex flex-wrap items-center justify-between gap-y-2",children:Object.keys(d).map((o,i)=>e.jsx("div",{className:x(p.option,p[`${o}Icon`],s===o&&p.active),onClick:()=>{C(o),E()}},i))}),s==="chromePlugin"&&e.jsx("div",{className:"mt-6 w-full",children:e.jsxs("div",{className:x("inline-flex w-full items-center justify-center gap-2 rounded-lg py-3","shrink-0 cursor-pointer bg-primary-600 text-white hover:bg-primary-600/75 hover:shadow-sm"),children:[e.jsx("div",{className:`relative h-4 w-4 ${p.pluginInstallIcon}`}),e.jsx("div",{className:"font-['Inter'] text-sm font-medium leading-tight text-white",onClick:I,children:n(`${a}.chromePlugin`,{ns:"appOverview"})})]})}),e.jsxs("div",{className:x("inline-flex w-full flex-col items-start justify-start rounded-lg border-[0.5px] border-components-panel-border bg-background-section","mt-6"),children:[e.jsxs("div",{className:"inline-flex items-center justify-start gap-2 self-stretch rounded-t-lg bg-background-section-burn py-1  pl-3 pr-1",children:[e.jsx("div",{className:"system-sm-medium shrink-0 grow text-text-secondary",children:n(`${a}.${s}`,{ns:"appOverview"})}),e.jsx(R,{popupContent:(l[s]?n(`${a}.copied`,{ns:"appOverview"}):n(`${a}.copy`,{ns:"appOverview"}))||"",children:e.jsx(O,{children:e.jsxs("div",{onClick:_,children:[l[s]&&e.jsx(k,{className:"h-4 w-4"}),!l[s]&&e.jsx(T,{className:"h-4 w-4"})]})})})]}),e.jsx("div",{className:"flex w-full items-start justify-start gap-2 overflow-x-auto p-3",children:e.jsx("div",{className:"shrink grow basis-0 font-mono text-[13px] leading-tight text-text-secondary",children:e.jsx("pre",{className:"select-text",children:d[s].getContent(c,u,((y=g.theme)==null?void 0:y.primaryColor)??"#1C64F2",w)})})})]})]})};export{K as E};
