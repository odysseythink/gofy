import{r as u,U as O,u as Te,j as e,c as R,Y as We,A as $,bC as $e,L as He}from"./index-C6MQE079.js";import{u as Fe}from"./store-BvNu8ZCy.js";import{a as Ve}from"./index-BnCt45Lk.js";import{S as Ke}from"./secret-key-button-45qFFLxR.js";import{a as Qe}from"./lib.client-jc6rCG4d.js";import{u as Je}from"./use-theme-ctHsEsh4.js";import{w as Ye}from"./secret-key-modal-BolZ7k5X.js";import{$ as De,a as Ze,w as en,e as nn}from"./use-resolve-button-type-KuxMC1Tz.js";import{Y as V,d as Y,y as K,j as Q,o as G,n as ee,G as H,K as J,a as F,V as Ie,a2 as Z,U as sn,l as I,W as B,D as P,a5 as rn,A as te}from"./portal-BKmWNQhF.js";import{f as dn}from"./transition-4KkSNCv7.js";import{f as Pe,s as ln}from"./active-element-history-CnAoIe1W.js";import"./index-C5I3Iaw-.js";import"./index-81ZfuQxj.js";import"./index-DwmcDQjG.js";import"./tslib.es6-CyuPUxRh.js";import"./index-Cy3FN14p.js";import"./index-CXcgbTpF.js";import"./index-CiU3yiP3.js";import"./index-B0HEP2-Y.js";import"./use-tab-direction-CyM8LT4Y.js";import"./app-context-BaKmc3Lv.js";import"./use-timestamp-ChQDKMtp.js";import"./dayjs.min-DdW-s3EH.js";import"./dayjs.min-B2iisHW2.js";import"./timezone-B1RYvtJa.js";import"./utc-7AxelFwI.js";import"./apps-A2MhT6Lj.js";import"./datasets-WJyPv7j-.js";import"./index-D2wPdbFB.js";import"./use-dataset-CTDi6rr2.js";import"./useMutation-GOX7jtru.js";import"./useInfiniteQuery-Cj2t0p1u.js";import"./use-base-5uJnQsuZ.js";import"./use-apps-rKAYTjYS.js";import"./XMarkIcon-C5uHLoJl.js";import"./PlusIcon-CqB_p-10.js";function cn({onFocus:i}){let[n,a]=u.useState(!0),o=dn();return n?O.createElement(Pe,{as:"button",type:"button",features:ln.Focusable,onFocus:h=>{h.preventDefault();let x,j=50;function _(){if(j--<=0){x&&cancelAnimationFrame(x);return}if(i()){if(cancelAnimationFrame(x),!o.current)return;a(!1);return}x=requestAnimationFrame(_)}x=requestAnimationFrame(_)}}):null}const Ae=u.createContext(null);function tn(){return{groups:new Map,get(i,n){var a;let o=this.groups.get(i);o||(o=new Map,this.groups.set(i,o));let h=(a=o.get(n))!=null?a:0;o.set(n,h+1);let x=Array.from(o.keys()).indexOf(n);function j(){let _=o.get(n);_>1?o.set(n,_-1):o.delete(n)}return[x,j]}}}function an({children:i}){let n=u.useRef(tn());return u.createElement(Ae.Provider,{value:n},i)}function Re(i){let n=u.useContext(Ae);if(!n)throw new Error("You must wrap your component in a <StableCollection>");let a=u.useId(),[o,h]=n.current.get(i,a);return u.useEffect(()=>h,[]),o}var on=(i=>(i[i.Forwards=0]="Forwards",i[i.Backwards=1]="Backwards",i))(on||{}),hn=(i=>(i[i.Less=-1]="Less",i[i.Equal=0]="Equal",i[i.Greater=1]="Greater",i))(hn||{}),xn=(i=>(i[i.SetSelectedIndex=0]="SetSelectedIndex",i[i.RegisterTab=1]="RegisterTab",i[i.UnregisterTab=2]="UnregisterTab",i[i.RegisterPanel=3]="RegisterPanel",i[i.UnregisterPanel=4]="UnregisterPanel",i))(xn||{});let jn={0(i,n){var a;let o=H(i.tabs,p=>p.current),h=H(i.panels,p=>p.current),x=o.filter(p=>{var g;return!((g=p.current)!=null&&g.hasAttribute("disabled"))}),j={...i,tabs:o,panels:h};if(n.index<0||n.index>o.length-1){let p=F(Math.sign(n.index-i.selectedIndex),{[-1]:()=>1,0:()=>F(Math.sign(n.index),{[-1]:()=>0,0:()=>0,1:()=>1}),1:()=>0});if(x.length===0)return j;let g=F(p,{0:()=>o.indexOf(x[0]),1:()=>o.indexOf(x[x.length-1])});return{...j,selectedIndex:g===-1?i.selectedIndex:g}}let _=o.slice(0,n.index),y=[...o.slice(n.index),..._].find(p=>x.includes(p));if(!y)return j;let m=(a=o.indexOf(y))!=null?a:i.selectedIndex;return m===-1&&(m=i.selectedIndex),{...j,selectedIndex:m}},1(i,n){if(i.tabs.includes(n.tab))return i;let a=i.tabs[i.selectedIndex],o=H([...i.tabs,n.tab],x=>x.current),h=i.selectedIndex;return i.info.current.isControlled||(h=o.indexOf(a),h===-1&&(h=i.selectedIndex)),{...i,tabs:o,selectedIndex:h}},2(i,n){return{...i,tabs:i.tabs.filter(a=>a!==n.tab)}},3(i,n){return i.panels.includes(n.panel)?i:{...i,panels:H([...i.panels,n.panel],a=>a.current)}},4(i,n){return{...i,panels:i.panels.filter(a=>a!==n.panel)}}},re=u.createContext(null);re.displayName="TabsDataContext";function U(i){let n=u.useContext(re);if(n===null){let a=new Error(`<${i} /> is missing a parent <Tab.Group /> component.`);throw Error.captureStackTrace&&Error.captureStackTrace(a,U),a}return n}let de=u.createContext(null);de.displayName="TabsActionsContext";function le(i){let n=u.useContext(de);if(n===null){let a=new Error(`<${i} /> is missing a parent <Tab.Group /> component.`);throw Error.captureStackTrace&&Error.captureStackTrace(a,le),a}return n}function un(i,n){return F(n.type,jn,i,n)}let pn="div";function _n(i,n){let{defaultIndex:a=0,vertical:o=!1,manual:h=!1,onChange:x,selectedIndex:j=null,..._}=i;const y=o?"vertical":"horizontal",m=h?"manual":"auto";let p=j!==null,g=Y({isControlled:p}),k=K(n),[b,f]=u.useReducer(un,{info:g,selectedIndex:j??a,tabs:[],panels:[]}),M=Q({selectedIndex:b.selectedIndex}),N=Y(x||(()=>{})),C=Y(b.tabs),w=u.useMemo(()=>({orientation:y,activation:m,...b}),[y,m,b]),L=G(v=>(f({type:1,tab:v}),()=>f({type:2,tab:v}))),X=G(v=>(f({type:3,panel:v}),()=>f({type:4,panel:v}))),T=G(v=>{D.current!==v&&N.current(v),p||f({type:0,index:v})}),D=Y(p?i.selectedIndex:b.selectedIndex),E=u.useMemo(()=>({registerTab:L,registerPanel:X,change:T}),[]);ee(()=>{f({type:0,index:j??a})},[j]),ee(()=>{if(D.current===void 0||b.tabs.length<=0)return;let v=H(b.tabs,S=>S.current);v.some((S,z)=>b.tabs[z]!==S)&&T(v.indexOf(b.tabs[D.current]))});let ne={ref:k},W=J();return O.createElement(an,null,O.createElement(de.Provider,{value:E},O.createElement(re.Provider,{value:w},w.tabs.length<=0&&O.createElement(cn,{onFocus:()=>{var v,S;for(let z of C.current)if(((v=z.current)==null?void 0:v.tabIndex)===0)return(S=z.current)==null||S.focus(),!0;return!1}}),W({ourProps:ne,theirProps:_,slot:M,defaultTag:pn,name:"Tabs"}))))}let mn="div";function gn(i,n){let{orientation:a,selectedIndex:o}=U("Tab.List"),h=K(n),x=Q({selectedIndex:o}),j=i,_={ref:h,role:"tablist","aria-orientation":a};return J()({ourProps:_,theirProps:j,slot:x,defaultTag:mn,name:"Tabs.List"})}let bn="button";function fn(i,n){var a,o;let h=u.useId(),{id:x=`headlessui-tabs-tab-${h}`,disabled:j=!1,autoFocus:_=!1,...y}=i,{orientation:m,activation:p,selectedIndex:g,tabs:k,panels:b}=U("Tab"),f=le("Tab"),M=U("Tab"),[N,C]=u.useState(null),w=u.useRef(null),L=K(w,n,C);ee(()=>f.registerTab(w),[f,w]);let X=Re("tabs"),T=k.indexOf(w);T===-1&&(T=X);let D=T===g,E=G(q=>{let A=q();if(A===Z.Success&&p==="auto"){let se=sn(w.current),ce=M.tabs.findIndex(Xe=>Xe.current===se);ce!==-1&&f.change(ce)}return A}),ne=G(q=>{let A=k.map(se=>se.current).filter(Boolean);if(q.key===I.Space||q.key===I.Enter){q.preventDefault(),q.stopPropagation(),f.change(T);return}switch(q.key){case I.Home:case I.PageUp:return q.preventDefault(),q.stopPropagation(),E(()=>B(A,P.First));case I.End:case I.PageDown:return q.preventDefault(),q.stopPropagation(),E(()=>B(A,P.Last))}if(E(()=>F(m,{vertical(){return q.key===I.ArrowUp?B(A,P.Previous|P.WrapAround):q.key===I.ArrowDown?B(A,P.Next|P.WrapAround):Z.Error},horizontal(){return q.key===I.ArrowLeft?B(A,P.Previous|P.WrapAround):q.key===I.ArrowRight?B(A,P.Next|P.WrapAround):Z.Error}}))===Z.Success)return q.preventDefault()}),W=u.useRef(!1),v=G(()=>{var q;W.current||(W.current=!0,(q=w.current)==null||q.focus({preventScroll:!0}),f.change(T),rn(()=>{W.current=!1}))}),S=G(q=>{q.preventDefault()}),{isFocusVisible:z,focusProps:Me}=De({autoFocus:_}),{isHovered:ze,hoverProps:Be}=Ze({isDisabled:j}),{pressed:Oe,pressProps:Ue}=en({disabled:j}),Ne=Q({selected:D,hover:ze,active:Oe,focus:z,autofocus:_,disabled:j}),Le=Ie({ref:L,onKeyDown:ne,onMouseDown:S,onClick:v,id:x,role:"tab",type:nn(i,N),"aria-controls":(o=(a=b[T])==null?void 0:a.current)==null?void 0:o.id,"aria-selected":D,tabIndex:D?0:-1,disabled:j||void 0,autoFocus:_},Me,Be,Ue);return J()({ourProps:Le,theirProps:y,slot:Ne,defaultTag:bn,name:"Tabs.Tab"})}let qn="div";function yn(i,n){let{selectedIndex:a}=U("Tab.Panels"),o=K(n),h=Q({selectedIndex:a}),x=i,j={ref:o};return J()({ourProps:j,theirProps:x,slot:h,defaultTag:qn,name:"Tabs.Panels"})}let vn="div",kn=te.RenderStrategy|te.Static;function wn(i,n){var a,o,h,x;let j=u.useId(),{id:_=`headlessui-tabs-panel-${j}`,tabIndex:y=0,...m}=i,{selectedIndex:p,tabs:g,panels:k}=U("Tab.Panel"),b=le("Tab.Panel"),f=u.useRef(null),M=K(f,n);ee(()=>b.registerPanel(f),[b,f]);let N=Re("panels"),C=k.indexOf(f);C===-1&&(C=N);let w=C===p,{isFocusVisible:L,focusProps:X}=De(),T=Q({selected:w,focus:L}),D=Ie({ref:M,id:_,role:"tabpanel","aria-labelledby":(o=(a=g[C])==null?void 0:a.current)==null?void 0:o.id,tabIndex:w?y:-1},X),E=J();return!w&&((h=m.unmount)==null||h)&&!((x=m.static)!=null&&x)?O.createElement(Pe,{"aria-hidden":"true",...D}):E({ourProps:D,theirProps:m,slot:T,defaultTag:vn,features:kn,visible:w,name:"Tabs.Panel"})}let Tn=V(fn),Ce=V(_n),Se=V(gn),Ee=V(yn),Ge=V(wn),Dn=Object.assign(Tn,{Group:Ce,List:Se,Panels:Ee,Panel:Ge});const In=({apiBaseUrl:i,appId:n})=>{const{t:a}=Te();return e.jsxs("div",{className:"flex flex-wrap items-center gap-y-2",children:[e.jsxs("div",{className:"mr-2 flex h-8 items-center rounded-lg border-[0.5px] border-components-input-border-active bg-components-input-bg-normal pl-1.5 pr-1 leading-5",children:[e.jsx("div",{className:"mr-0.5 h-5 shrink-0 rounded-md border border-divider-subtle px-1.5 text-[11px] text-text-tertiary",children:a("apiServer",{ns:"appApi"})}),e.jsx("div",{className:"w-fit truncate px-1 text-[13px] font-medium text-text-secondary sm:w-[248px]",children:i}),e.jsx("div",{className:"mx-1 h-[14px] w-[1px] bg-divider-regular"}),e.jsx(Ve,{content:i})]}),e.jsx("div",{className:"mr-2 flex h-8 items-center rounded-lg border-[0.5px] border-[#D1FADF] bg-[#ECFDF3] px-3 text-xs font-semibold text-[#039855]",children:a("ok",{ns:"appApi"})}),e.jsx(Ke,{className:"!h-8 shrink-0",appId:n})]})},Pn=80,ae=".overflow-auto",ie=i=>i.replace("#",""),An=()=>{const i=document.querySelector("article");return i?Array.from(i.querySelectorAll("h2")).map(n=>{const a=n.querySelector("a");return a?{href:a.getAttribute("href")||"",text:a.textContent||""}:null}).filter(n=>n!==null):[]},Rn=({appDetail:i,locale:n})=>{const[a,o]=u.useState([]),[h,x]=u.useState(()=>typeof window>"u"?!1:window.matchMedia("(min-width: 1280px)").matches),[j,_]=u.useState("");u.useEffect(()=>{const m=setTimeout(()=>{const p=An();o(p),p.length>0&&_(ie(p[0].href))},0);return()=>clearTimeout(m)},[i,n]),u.useEffect(()=>{const m=document.querySelector(ae);if(!m||a.length===0)return;const p=()=>{let g="";for(const k of a){const b=ie(k.href),f=document.getElementById(b);f&&f.getBoundingClientRect().top<=window.innerHeight/2&&(g=b)}g&&g!==j&&_(g)};return m.addEventListener("scroll",p),()=>m.removeEventListener("scroll",p)},[a,j]);const y=u.useCallback((m,p)=>{m.preventDefault();const g=ie(p.href),k=document.getElementById(g);if(!k)return;const b=document.querySelector(ae);b&&b.scrollTo({top:k.offsetTop-Pn,behavior:"smooth"})},[]);return{toc:a,isTocExpanded:h,setIsTocExpanded:x,activeSection:j,handleTocClick:y}},Cn={medium:"rounded-lg px-1.5 ring-1 ring-inset"},Sn={emerald:{small:"text-emerald-500 dark:text-emerald-400",medium:"ring-emerald-300 dark:ring-emerald-400/30 bg-emerald-400/10 text-emerald-500 dark:text-emerald-400"},sky:{small:"text-sky-500",medium:"ring-sky-300 bg-sky-400/10 text-sky-500 dark:ring-sky-400/30 dark:bg-sky-400/10 dark:text-sky-400"},amber:{small:"text-amber-500",medium:"ring-amber-300 bg-amber-400/10 text-amber-500 dark:ring-amber-400/30 dark:bg-amber-400/10 dark:text-amber-400"},rose:{small:"text-red-500 dark:text-rose-500",medium:"ring-rose-200 bg-rose-50 text-red-500 dark:ring-rose-500/20 dark:bg-rose-400/10 dark:text-rose-400"},zinc:{small:"text-zinc-400 dark:text-zinc-500",medium:"ring-zinc-200 bg-zinc-50 text-zinc-500 dark:ring-zinc-500/20 dark:bg-zinc-400/10 dark:text-zinc-400"}},En={get:"emerald",post:"sky",put:"amber",delete:"rose"};function Gn({children:i,variant:n="medium",color:a=En[i.toLowerCase()]??"emerald"}){return e.jsx("span",{className:R("font-mono text-[0.625rem] font-semibold leading-6",Cn[n],Sn[a][n]),children:i})}function Mn(i){return e.jsxs("svg",{viewBox:"0 0 20 20","aria-hidden":"true",...i,children:[e.jsx("path",{strokeWidth:"0",d:"M5.5 13.5v-5a2 2 0 0 1 2-2l.447-.894A2 2 0 0 1 9.737 4.5h.527a2 2 0 0 1 1.789 1.106l.447.894a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-5a2 2 0 0 1-2-2Z"}),e.jsx("path",{fill:"none",strokeLinejoin:"round",d:"M12.5 6.5a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-5a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2m5 0-.447-.894a2 2 0 0 0-1.79-1.106h-.527a2 2 0 0 0-1.789 1.106L7.5 6.5m5 0-1 1h-3l-1-1"})]})}function zn({code:i}){const[n,a]=u.useState(0),o=n>0;return u.useEffect(()=>{if(n>0){const h=setTimeout(()=>a(0),1e3);return()=>{clearTimeout(h)}}},[n]),e.jsxs("button",{type:"button",className:R("group/button absolute right-4 top-1.5 overflow-hidden rounded-full py-1 pl-2 pr-3 text-2xs font-medium opacity-0 backdrop-blur transition focus:opacity-100 group-hover:opacity-100",o?"bg-emerald-400/10 ring-1 ring-inset ring-emerald-400/20":"hover:bg-white/7.5 dark:bg-white/2.5 bg-white/5 dark:hover:bg-white/5"),onClick:()=>{Ye(i).then(()=>{a(h=>h+1)})},children:[e.jsxs("span",{"aria-hidden":o,className:R("pointer-events-none flex items-center gap-0.5 text-zinc-400 transition duration-300",o&&"-translate-y-1.5 opacity-0"),children:[e.jsx(Mn,{className:"h-5 w-5 fill-zinc-500/20 stroke-zinc-500 transition-colors group-hover/button:stroke-zinc-400"}),"Copy"]}),e.jsx("span",{"aria-hidden":!o,className:R("pointer-events-none absolute inset-0 flex items-center justify-center text-emerald-400 transition duration-300",!o&&"translate-y-1.5 opacity-0"),children:"Copied!"})]})}function Bn({tag:i,label:n}){return!i&&!n?null:e.jsxs("div",{className:"border-b-white/7.5 bg-white/2.5 dark:bg-white/1 flex h-9 items-center gap-2 border-y border-t-transparent bg-zinc-900 px-4 dark:border-b-white/5",children:[i&&e.jsx("div",{className:"dark flex",children:e.jsx(Gn,{variant:"small",children:i})}),i&&n&&e.jsx("span",{className:"h-0.5 w-0.5 rounded-full bg-zinc-500"}),n&&e.jsx("span",{className:"font-mono text-xs text-zinc-400",children:n})]})}function oe({tag:i,label:n,children:a,targetCode:o}){const h=u.Children.toArray(a)[0];return e.jsxs("div",{className:"dark:bg-white/2.5 group",children:[e.jsx(Bn,{tag:i,label:n}),e.jsxs("div",{className:"relative",children:[e.jsx("pre",{className:"overflow-x-auto p-4 text-xs text-white",children:o!=null&&o.code?e.jsx("code",{children:o==null?void 0:o.code}):h}),e.jsx(zn,{code:(o==null?void 0:o.code)??h.props.children.props.children})]})]})}function On({title:i,tabTitles:n,selectedIndex:a}){const o=((n==null?void 0:n.length)??0)>1;return e.jsxs("div",{className:"flex min-h-[calc(theme(spacing.12)+1px)] flex-wrap items-start gap-x-4 border-b border-zinc-700 bg-zinc-800 px-4 dark:border-zinc-800 dark:bg-transparent",children:[i&&e.jsx("h3",{className:"mr-auto pt-3 text-xs font-semibold text-white",children:i}),o&&e.jsx(Se,{className:"-mb-px flex gap-4 text-xs font-medium",children:n.map((h,x)=>e.jsx(Dn,{className:R("border-b py-3 transition focus:[&:not(:focus-visible)]:outline-none",x===a?"border-emerald-500 text-emerald-400":"border-transparent text-zinc-400 hover:text-zinc-300"),children:h},x))})]})}function Un({children:i,targetCode:n,...a}){return((n==null?void 0:n.length)??0)>1?e.jsx(Ee,{children:n.map((o,h)=>e.jsx(Ge,{children:e.jsx(oe,{...a,targetCode:o})},o.title||o.tag||h))}):e.jsx(oe,{...a,targetCode:n==null?void 0:n[0],children:i})}function Nn(){const i=u.useRef(null),n=u.useRef(null);return u.useEffect(()=>()=>{window.cancelAnimationFrame(n.current)},[]),{positionRef:i,preventLayoutShift(a){const o=i.current.getBoundingClientRect().top;a(),n.current=window.requestAnimationFrame(()=>{const h=i.current.getBoundingClientRect().top;window.scrollBy(0,h-o)})}}}function Ln(i){const[n,a]=u.useState([]),[o,h]=u.useState(0),x=[...i||[]].sort((p,g)=>n.indexOf(g)-n.indexOf(p))[0],j=(i==null?void 0:i.indexOf(x))||0,_=j===-1?o:j;_!==o&&h(_);const{positionRef:y,preventLayoutShift:m}=Nn();return{as:"div",ref:y,selectedIndex:o,onChange:p=>{m(()=>a(i[p]))}}}const Xn=u.createContext(!1);function s({children:i,title:n,targetCode:a,...o}){const h=typeof a=="string"?[{code:a}]:a,x=(h==null?void 0:h.map(({title:g})=>g||"Code"))||[],j=Ln(x),_=x.length>1,y=_?Ce:"div",m=_?j:{},p=_?{selectedIndex:j.selectedIndex,tabTitles:x}:{};return e.jsx(Xn.Provider,{value:!0,children:e.jsxs(y,{...m,className:"not-prose my-6 overflow-hidden rounded-2xl bg-zinc-900 shadow-md dark:ring-1 dark:ring-white/10",children:[e.jsx(On,{title:n,...p}),e.jsx(Un,{...o,targetCode:h,children:i})]})})}const l=function({url:n,method:a,title:o,name:h}){let x="";switch(a){case"PUT":x="ring-amber-300 bg-amber-400/10 text-amber-500 dark:ring-amber-400/30 dark:bg-amber-400/10 dark:text-amber-400";break;case"DELETE":x="ring-rose-200 bg-rose-50 text-red-500 dark:ring-rose-500/20 dark:bg-rose-400/10 dark:text-rose-400";break;case"POST":x="ring-sky-300 bg-sky-400/10 text-sky-500 dark:ring-sky-400/30 dark:bg-sky-400/10 dark:text-sky-400";break;case"PATCH":x="ring-violet-300 bg-violet-400/10 text-violet-500 dark:ring-violet-400/30 dark:bg-violet-400/10 dark:text-violet-400";break;default:x="ring-emerald-300 dark:ring-emerald-400/30 bg-emerald-400/10 text-emerald-500 dark:text-emerald-400";break}return e.jsxs(e.Fragment,{children:[e.jsx("span",{id:h==null?void 0:h.replace(/^#/,""),className:"relative -top-28"}),e.jsxs("div",{className:"flex items-center gap-x-3",children:[e.jsx("span",{className:`rounded-lg px-1.5 font-mono text-[0.625rem] font-semibold leading-6 ring-1 ring-inset ${x}`,children:a}),e.jsx("span",{className:"font-mono text-xs text-zinc-400",children:n})]}),e.jsx("h2",{className:"mt-2 scroll-mt-32",children:e.jsx("a",{href:h,className:"group text-inherit no-underline hover:text-inherit",children:o})})]})};function c({children:i}){return e.jsx("div",{className:"grid grid-cols-1 items-start gap-x-16 gap-y-10 xl:!max-w-none xl:grid-cols-2",children:i})}function r({children:i,sticky:n=!1}){return e.jsx("div",{className:R("[&>:first-child]:mt-0 [&>:last-child]:mb-0",n&&"xl:sticky xl:top-24"),children:i})}function t({children:i}){return e.jsx("div",{className:"my-6",children:e.jsx("ul",{role:"list",className:"m-0 max-w-[calc(theme(maxWidth.lg)-theme(spacing.8))] list-none divide-y divide-zinc-900/5 p-0 dark:divide-white/5",children:i})})}function d({name:i,type:n,children:a}){return e.jsx("li",{className:"m-0 px-0 py-4 first:pt-0 last:pb-0",children:e.jsxs("dl",{className:"m-0 flex flex-wrap items-center gap-x-3 gap-y-2",children:[e.jsx("dt",{className:"sr-only",children:"Name"}),e.jsx("dd",{children:e.jsx("code",{children:i})}),e.jsx("dt",{className:"sr-only",children:"Type"}),e.jsx("dd",{className:"font-mono text-xs text-zinc-400 dark:text-zinc-500",children:n}),e.jsx("dt",{className:"sr-only",children:"Description"}),e.jsx("dd",{className:"w-full flex-none [&>:first-child]:mt-0 [&>:last-child]:mb-0",children:a})]})})}function he(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"Completion App API"}),`
`,e.jsx(n.p,{children:"The text generation application offers non-session support and is ideal for translation, article writing, summarization AI, and more."}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"Base URL"}),e.jsx(s,{title:"Code",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"Authentication"}),e.jsxs(n.p,{children:["The Service API uses ",e.jsx(n.code,{children:"API-Key"}),` authentication.
`,e.jsx("i",{children:e.jsx(n.strong,{children:"Strongly recommend storing your API Key on the server-side, not shared or stored on the client-side, to avoid possible API-Key leakage that can lead to serious consequences."})})]}),e.jsxs(n.p,{children:["For all API requests, include your API Key in the ",e.jsx(n.code,{children:"Authorization"})," HTTP Header, as shown below:"]}),e.jsx(s,{title:"Code",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/completion-messages",method:"POST",title:"Create Completion Message",name:"#Create-Completion-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Send a request to the text generation application."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsxs(d,{name:"inputs",type:"object",children:[e.jsxs(n.p,{children:[`Allows the entry of various variable values defined by the App.
The `,e.jsx(n.code,{children:"inputs"}),` parameter contains multiple key/value pairs, with each key corresponding to a specific variable and each value being the specific value for that variable.
The text generation application requires at least one key/value pair to be inputted.`]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"}),` (string) Required
The input text, the content to be processed.`]}),`
`]})]},"inputs"),e.jsxs(d,{name:"response_mode",type:"string",children:[e.jsx(n.p,{children:"The mode of response return, supporting:"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," Streaming mode (recommended), implements a typewriter-like output through SSE (",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` Blocking mode, returns result after execution is complete. (Requests may be interrupted if the process is long)
`,e.jsx("i",{children:"Due to Cloudflare restrictions, the request will be interrupted without a return after 100 seconds."})]}),`
`]})]},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user, convenient for retrieval and statistics.
The rules are defined by the developer and need to ensure that the user identifier is unique within the application. The Service API does not share conversations created by the WebApp.`})},"user"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"File list, suitable for inputting files combined with text understanding and answering questions, available only when the model supports Vision/Video capability."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) Supported type:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," Supported types include: 'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," Supported types include: 'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," Supported types include: 'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," Supported types include: 'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," Supported types include: other file types"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) Transfer method:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": File URL."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": Upload file."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," File URL. (Only when transfer method is ",e.jsx(n.code,{children:"remote_url"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," Upload file ID. (Only when transfer method is ",e.jsx(n.code,{children:"local_file"}),")."]}),`
`]})]},"files")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.p,{children:["When ",e.jsx(n.code,{children:"response_mode"})," is ",e.jsx(n.code,{children:"blocking"}),`, return a CompletionResponse object.
When `,e.jsx(n.code,{children:"response_mode"})," is ",e.jsx(n.code,{children:"streaming"}),", return a ChunkCompletionResponse stream."]}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["Returns the complete App result, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"application/json"}),"."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) App mode, fixed as ",e.jsx(n.code,{children:"chat"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) Complete response content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) Metadata",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) Model usage information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) Citation and Attribution List"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Message creation timestamp, e.g., 1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["Returns the stream chunks outputted by the App, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"text/event-stream"}),`.
Each streaming chunk starts with `,e.jsx(n.code,{children:"data:"}),", separated by two newline characters ",e.jsx(n.code,{children:"\\n\\n"}),", as shown below:"]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsxs(n.p,{children:["The structure of the streaming chunks varies depending on the ",e.jsx(n.code,{children:"event"}),":"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLM returns text chunk event, i.e., the complete text is output in a chunked fashion.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLM returned text chunk content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," Message end event, receiving this event means streaming has ended.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) Metadata",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) Model usage information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) Citation and Attribution List"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS audio stream event, that is, speech synthesis output. The content is an audio block in Mp3 format, encoded as a base64 string. When playing, simply decode the base64 and feed it into the player. (This message is available only when auto-play is enabled)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the stop response interface below"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The audio after speech synthesis, encoded in base64 text content, when playing, simply decode the base64 and feed it into the player"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g.: 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS audio stream end event, receiving this event indicates the end of the audio stream.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the stop response interface below"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The end event has no audio, so this is an empty string"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g.: 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` Message content replacement event.
When output content moderation is enabled, if the content is flagged, then the message content will be replaced with a preset reply through this event.`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) Replacement content (directly replaces all LLM reply text)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
Exceptions that occur during the streaming process will be output in the form of stream events, and reception of an error event will end the stream.`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTP status code"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) Error code"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) Error message"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," Ping event every 10 seconds to keep the connection alive."]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404, Conversation does not exists"}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", abnormal parameter input"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", App configuration unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", no available model credential configuration"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", model invocation quota insufficient"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", current model unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"completion_request_error"}),", text generation failed"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/completion-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/completion-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": {"query": "Hello, world!"},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(n.h3,{children:"Blocking Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "event": "message",
  "message_id": "9da23599-e713-473b-982c-4328d4f5c78a",
  "mode": "completion",
  "answer": "Hello World!...",
  "metadata": {
      "usage": {
          "prompt_tokens": 1033,
          "prompt_unit_price": "0.001",
          "prompt_price_unit": "0.001",
          "prompt_price": "0.0010330",
          "completion_tokens": 128,
          "completion_unit_price": "0.002",
          "completion_price_unit": "0.001",
          "completion_price": "0.0002560",
          "total_tokens": 1161,
          "total_price": "0.0012890",
          "currency": "USD",
          "latency": 0.7682376249867957
      }
  },
  "created_at": 1705407629
}
`})})}),e.jsx(n.h3,{children:"Streaming Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " I", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": "'m", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " glad", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " to", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " meet", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " you", "created_at": 1679586595}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"File Upload",name:"#file-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:[`Upload a file (currently only images are supported) for use when sending messages, enabling multimodal understanding of images and text.
Supports png, jpg, jpeg, webp, gif formats.
`,e.jsx("i",{children:"Uploaded files are for use by the current end-user only."})]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["This interface requires a ",e.jsx(n.code,{children:"multipart/form-data"})," request."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file"}),` (File) Required
The file to be uploaded.`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, defined by the developer's rules, must be unique within the application. The Service API does not share conversations created by the WebApp.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"After a successful upload, the server will return the file's ID and related information."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) File name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) File size (bytes)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) File extension"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) File mime-type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) End-user ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"no_file_uploaded"}),", a file must be provided"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"too_many_files"}),", currently only one file is accepted"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_preview"}),", the file does not support preview"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_estimate"}),", the file does not support estimation"]}),`
`,e.jsxs(n.li,{children:["413, ",e.jsx(n.code,{children:"file_too_large"}),", the file is too large"]}),`
`,e.jsxs(n.li,{children:["415, ",e.jsx(n.code,{children:"unsupported_file_type"}),", unsupported extension, currently only document files are accepted"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_connection_failed"}),", unable to connect to S3 service"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_permission_denied"}),", no permission to upload files to S3"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_file_too_large"}),", file exceeds S3 size limit"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"Get End User",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve an end user by ID."}),e.jsxs(n.p,{children:["This is useful when other APIs return an end-user ID (e.g. ",e.jsx(n.code,{children:"created_by"})," from File Upload)."]}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) Required
End user ID.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns an EndUser object."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) Tenant ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) App ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) End user type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) External user ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) Whether anonymous"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) Session ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 datetime"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 datetime"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"end_user_not_found"}),", end user not found"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"File Preview",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Preview or download uploaded files. This endpoint allows you to access files that have been previously uploaded via the File Upload API."}),e.jsx("i",{children:"Files can only be accessed if they belong to messages within the requesting application."}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) Required
The unique identifier of the file to preview, obtained from the File Upload API response.`]}),`
`]}),e.jsx(n.h3,{children:"Query Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) Optional
Whether to force download the file as an attachment. Default is `,e.jsx(n.code,{children:"false"})," (preview in browser)."]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns the file content with appropriate headers for browser display or download."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," Set based on file mime type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," File size in bytes (if available)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"}),' Set to "attachment" if ',e.jsx(n.code,{children:"as_attachment=true"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," Caching headers for performance"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' Set to "bytes" for audio/video files']}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", abnormal parameter input"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", file access denied or file does not belong to current application"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", file not found or has been deleted"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Download as Attachment"}),e.jsx(s,{title:`Download
Request`,tag:"GET",label:"/files/:file_id/preview?as_attachment=true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"Response Headers Example"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"Download Response Headers"}),e.jsx(s,{title:"Download Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/completion-messages/:task_id/stop",method:"POST",title:"Stop Generate",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Only supported in streaming mode."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"}),` (string) Task ID, can be obtained from the streaming chunk return
Request Body`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, used to define the identity of the end-user, must be consistent with the user passed in the send message interface. The Service API does not share conversations created by the WebApp.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/completion-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/completion-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{ "user": "abc-123"}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"Message Feedback",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"End-users can provide feedback messages, facilitating application developers to optimize expected outputs."}),e.jsx(n.h3,{children:"Path"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"Message ID"})},"message_id")}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsxs(n.p,{children:["Upvote as ",e.jsx(n.code,{children:"like"}),", downvote as ",e.jsx(n.code,{children:"dislike"}),", revoke upvote as ",e.jsx(n.code,{children:"null"})]})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"User identifier, defined by the developer's rules, must be unique within the application."})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"The specific content of message feedback."})},"content")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"Get feedbacks of application",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Get application's feedbacks."}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（optional）pagination，default：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（optional） records per page default：20"})},"limit")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (List) return apps feedback list."]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`  {
      "data": [
          {
              "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
              "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
              "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
              "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
              "rating": "like",
              "content": "message feedback information-3",
              "from_source": "user",
              "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
              "from_account_id": null,
              "created_at": "2025-04-24T09:24:38",
              "updated_at": "2025-04-24T09:24:38"
          }
      ]
  }
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"Text to Audio",name:"#audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Text to speech."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"For text messages generated by Gofy, simply pass the generated message-id directly. The backend will use the message-id to look up the corresponding content and synthesize the voice information directly. If both message_id and text are provided simultaneously, the message_id is given priority."})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"Speech generated content."})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the app."})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/text-to-audio",targetCode:`curl -o text-to-audio.mp3 -X POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290",
  "text": "Hello Gofy",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"Get Application Basic Information",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get basic information about this application"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) application name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) application description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) application tags"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) application mode"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"author_name"})," (string) author name"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"Get Application Parameters Information",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used at the start of entering the page to obtain information such as features, input parameter names, types, and default values."}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) Opening statement"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) List of suggested questions for the opening"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) Suggest questions after enabling the answer.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) Speech to text",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) Citation and Attribution",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) Annotation reply",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) User input form configuration",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) Text input control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) Paragraph text input control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) Dropdown control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) Option values"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) File upload configuration",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) Document settings
Currently only supports document types: `,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Document number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) Image settings
Currently only supports image types: `,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Image number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) Audio settings
Currently only supports audio types: `,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Audio number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) Video settings
Currently only supports video types: `,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Video number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) Custom settings",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Custom number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) System parameters",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) Document upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) Image file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) Audio file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) Video file upload size limit (MB)"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "opening_statement": "Hello!",
  "suggested_questions_after_answer": {
      "enabled": true
  },
  "speech_to_text": {
      "enabled": true
  },
  "retriever_resource": {
      "enabled": true
  },
  "annotation_reply": {
      "enabled": true
  },
  "user_input_form": [
      {
          "paragraph": {
              "label": "Query",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"Get Application WebApp Settings",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get the WebApp settings of the application."}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) Chat color theme, in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) Whether the chat color theme is inverted"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) Icon type, ",e.jsx(n.code,{children:"emoji"})," - emoji, ",e.jsx(n.code,{children:"image"})," - picture"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) Icon. If it's ",e.jsx(n.code,{children:"emoji"})," type, it's an emoji symbol; if it's ",e.jsx(n.code,{children:"image"})," type, it's an image URL."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) Background color in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) Icon URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) Description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) Copyright information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) Privacy policy link"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) Custom disclaimer"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) Default language"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) Whether to show workflow details"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) Whether to replace 🤖 in chat with the WebApp icon"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{})]})}function Wn(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(he,{...i})}):he(i)}function xe(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"Completion アプリ API"}),`
`,e.jsx(n.p,{children:"テキスト生成アプリケーションはセッションレスをサポートし、翻訳、記事作成、要約 AI 等に最適です。"}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"ベース URL"}),e.jsx(s,{title:"コード",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"認証"}),e.jsxs(n.p,{children:["サービス API は ",e.jsx(n.code,{children:"API-Key"}),` 認証を使用します。
`,e.jsx("i",{children:e.jsx(n.strong,{children:"API キーの漏洩による重大な結果を避けるため、API キーはサーバーサイドに保存し、クライアントサイドでは共有や保存しないことを強く推奨します。"})})]}),e.jsxs(n.p,{children:["すべての API リクエストで、以下のように ",e.jsx(n.code,{children:"Authorization"})," HTTP ヘッダーに API キーを含めてください："]}),e.jsx(s,{title:"コード",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/completion-messages",method:"POST",title:"完了メッセージの作成",name:"#Create-Completion-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"テキスト生成アプリケーションにリクエストを送信します。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsxs(d,{name:"inputs",type:"object",children:[e.jsxs(n.p,{children:[`アプリで定義された各種変数値を入力できます。
`,e.jsx(n.code,{children:"inputs"}),`パラメータには複数のキー/値ペアが含まれ、各キーは特定の変数に対応し、各値はその変数の具体的な値となります。
テキスト生成アプリケーションでは、少なくとも1つのキー/値ペアの入力が必要です。`]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"}),` (string) 必須
入力テキスト、処理される内容。`]}),`
`]})]},"inputs"),e.jsxs(d,{name:"response_mode",type:"string",children:[e.jsx(n.p,{children:"レスポンス返却モード、以下をサポート："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," ストリーミングモード（推奨）、SSE（",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"}),"）によるタイプライター風の出力を実装。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` ブロッキングモード、実行完了後に結果を返却。（処理が長い場合はリクエストが中断される可能性があります）
`,e.jsx("i",{children:"Cloudflareの制限により、100秒後に返却なしで中断されます。"})]}),`
`]})]},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーの身元を定義し、取得や統計に使用します。
アプリケーション内で開発者が一意に定義する必要があります。`})},"user"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"ファイルリスト、モデルが Vision/Video 機能をサポートしている場合に限り、ファイルをテキスト理解および質問応答に組み合わせて入力するのに適しています。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) サポートされるタイプ：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," サポートされるタイプには以下が含まれます：'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," サポートされるタイプには以下が含まれます：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," サポートされるタイプには以下が含まれます：'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," サポートされるタイプには以下が含まれます：'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," サポートされるタイプには以下が含まれます：その他のファイルタイプ"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) 転送方法:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": ファイルのURL。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": ファイルをアップロード。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," ファイルのURL。（転送方法が ",e.jsx(n.code,{children:"remote_url"})," の場合のみ）。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," アップロードされたファイルID。（転送方法が ",e.jsx(n.code,{children:"local_file"})," の場合のみ）。"]}),`
`]})]},"files")]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.p,{children:[e.jsx(n.code,{children:"response_mode"}),"が",e.jsx(n.code,{children:"blocking"}),`の場合、CompletionResponseオブジェクトを返却します。
`,e.jsx(n.code,{children:"response_mode"}),"が",e.jsx(n.code,{children:"streaming"}),"の場合、ChunkCompletionResponseストリームを返却します。"]}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["アプリの完全な結果を返却、",e.jsx(n.code,{children:"Content-Type"}),"は",e.jsx(n.code,{children:"application/json"}),"です。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) アプリモード、固定で",e.jsx(n.code,{children:"chat"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 完全な応答内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) モデル使用情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用と帰属のリスト"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) メッセージ作成タイムスタンプ、例：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["アプリが出力するストリームチャンクを返却、",e.jsx(n.code,{children:"Content-Type"}),"は",e.jsx(n.code,{children:"text/event-stream"}),`です。
各ストリーミングチャンクは`,e.jsx(n.code,{children:"data:"}),"で始まり、2つの改行文字",e.jsx(n.code,{children:"\\n\\n"}),"で区切られます："]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsxs(n.p,{children:["ストリーミングチャンクの構造は",e.jsx(n.code,{children:"event"}),"によって異なります："]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLMがテキストチャンクを返すイベント、つまり完全なテキストがチャンク形式で出力されます。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエストの追跡と以下の生成停止APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLMが返したテキストチャンクの内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," メッセージ終了イベント、このイベントを受信するとストリーミングが終了したことを意味します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエストの追跡と以下の生成停止APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) モデル使用情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用と帰属のリスト"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS音声ストリームイベント、つまり音声合成出力。内容はMp3形式の音声ブロックで、base64文字列としてエンコードされています。再生時は単にbase64をデコードしてプレーヤーに供給するだけです。（このメッセージは自動再生が有効な場合のみ利用可能）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエストの追跡と以下の応答停止インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 音声合成後の音声、base64テキストコンテンツとしてエンコード、再生時は単にbase64をデコードしてプレーヤーに供給"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS音声ストリーム終了イベント、このイベントを受信すると音声ストリームが終了したことを示します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエストの追跡と以下の応答停止インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 終了イベントには音声がないため、空文字列"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` メッセージ内容置換イベント。
出力内容のモデレーションが有効な場合、コンテンツがフラグ付けされると、このイベントを通じてメッセージ内容が事前設定された返信に置き換えられます。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエストの追跡と以下の生成停止APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 置換内容（LLMの返信テキストすべてを直接置換）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
ストリーミング処理中に発生した例外は、ストリームイベントの形式で出力され、エラーイベントを受信するとストリームが終了します。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエストの追跡と以下の生成停止APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTPステータスコード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) エラーコード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) エラーメッセージ"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 接続を維持するため10秒ごとのPingイベント。"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404, 会話が存在しません"}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", パラメータ入力異常"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", アプリ設定が利用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", 利用可能なモデル認証情報設定がありません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", モデル呼び出しクォータ不足"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", 現在のモデルは利用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"completion_request_error"}),", テキスト生成に失敗しました"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/completion-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/completion-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": {"query": "Hello, world!"},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(n.h3,{children:"ブロッキングモード"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "event": "message",
  "message_id": "9da23599-e713-473b-982c-4328d4f5c78a",
  "mode": "completion",
  "answer": "Hello World!...",
  "metadata": {
      "usage": {
          "prompt_tokens": 1033,
          "prompt_unit_price": "0.001",
          "prompt_price_unit": "0.001",
          "prompt_price": "0.0010330",
          "completion_tokens": 128,
          "completion_unit_price": "0.002",
          "completion_price_unit": "0.001",
          "completion_price": "0.0002560",
          "total_tokens": 1161,
          "total_price": "0.0012890",
          "currency": "USD",
          "latency": 0.7682376249867957
      }
  },
  "created_at": 1705407629
}
`})})}),e.jsx(n.h3,{children:"ストリーミングモード"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " I", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": "'m", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " glad", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " to", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " meet", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " you", "created_at": 1679586595}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"ファイルアップロード",name:"#file-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:[`メッセージ送信時に使用するファイル（現在は画像のみ対応）をアップロードし、画像とテキストのマルチモーダルな理解を可能にします。
png、jpg、jpeg、webp、gif 形式に対応しています。
`,e.jsx("i",{children:"アップロードされたファイルは、現在のエンドユーザーのみが使用できます。"})]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.p,{children:["このインターフェースは",e.jsx(n.code,{children:"multipart/form-data"}),"リクエストが必要です。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file"}),` (File) 必須
アップロードするファイル。`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
開発者のルールで定義されたユーザー識別子。アプリケーション内で一意である必要があります。サービス API は WebApp によって作成された会話を共有しません。`]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"アップロードが成功すると、サーバーはファイルの ID と関連情報を返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) ファイル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) ファイルサイズ（バイト）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) ファイル拡張子"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) ファイルの MIME タイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) エンドユーザーID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"no_file_uploaded"}),", ファイルを提供する必要があります"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"too_many_files"}),", 現在は 1 つのファイルのみ受け付けています"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_preview"}),", ファイルがプレビューに対応していません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_estimate"}),", ファイルが推定に対応していません"]}),`
`,e.jsxs(n.li,{children:["413, ",e.jsx(n.code,{children:"file_too_large"}),", ファイルが大きすぎます"]}),`
`,e.jsxs(n.li,{children:["415, ",e.jsx(n.code,{children:"unsupported_file_type"}),", サポートされていない拡張子です。現在はドキュメントファイルのみ受け付けています"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_connection_failed"}),", S3 サービスに接続できません"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_permission_denied"}),", S3 へのファイルアップロード権限がありません"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_file_too_large"}),", ファイルが S3 のサイズ制限を超えています"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(n.h3,{children:"レスポンス例"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"エンドユーザーを取得",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"エンドユーザー ID からエンドユーザー情報を取得します。"}),e.jsxs(n.p,{children:["他の API がエンドユーザー ID（例：ファイルアップロードの ",e.jsx(n.code,{children:"created_by"}),"）を返す場合に利用できます。"]}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) 必須
エンドユーザー ID。`]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"EndUser オブジェクトを返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) テナント ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) アプリ ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) エンドユーザー種別"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) 外部ユーザー ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) 匿名ユーザーかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) セッション ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 日時"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 日時"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"end_user_not_found"}),", エンドユーザーが見つかりません"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"レスポンス例"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"ファイルプレビュー",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アップロードされたファイルをプレビューまたはダウンロードします。このエンドポイントを使用すると、以前にファイルアップロード API でアップロードされたファイルにアクセスできます。"}),e.jsx("i",{children:"ファイルは、リクエストしているアプリケーションのメッセージ範囲内にある場合のみアクセス可能です。"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) 必須
プレビューするファイルの一意識別子。ファイルアップロード API レスポンスから取得します。`]}),`
`]}),e.jsx(n.h3,{children:"クエリパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) オプション
ファイルを添付ファイルとして強制ダウンロードするかどうか。デフォルトは `,e.jsx(n.code,{children:"false"}),"（ブラウザでプレビュー）。"]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"ブラウザ表示またはダウンロード用の適切なヘッダー付きでファイル内容を返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," ファイル MIME タイプに基づいて設定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," ファイルサイズ（バイト、利用可能な場合）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"})," ",e.jsx(n.code,{children:"as_attachment=true"}),' の場合は "attachment" に設定']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," パフォーマンス向上のためのキャッシュヘッダー"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' 音声/動画ファイルの場合は "bytes" に設定']}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", パラメータ入力異常"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", ファイルアクセス拒否またはファイルが現在のアプリケーションに属していません"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", ファイルが見つからないか削除されています"]}),`
`,e.jsx(n.li,{children:"500, サーバー内部エラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"添付ファイルとしてダウンロード"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"レスポンスヘッダー例"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"ファイルダウンロードレスポンスヘッダー"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/completion-messages/:task_id/stop",method:"POST",title:"生成の停止",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ストリーミングモードでのみサポートされています。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"}),` (string) タスク ID、ストリーミングチャンクの返信から取得可能
リクエストボディ`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子。エンドユーザーの身元を定義するために使用され、メッセージ送信インターフェースで渡されたユーザーと一致する必要があります。サービス API は WebApp によって作成された会話を共有しません。`]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"Request",tag:"POST",label:"/completion-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/completion-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{ "user": "abc-123"}'`}),e.jsx(n.h3,{children:"レスポンス例"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"メッセージフィードバック",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"エンドユーザーはフィードバックメッセージを提供でき、アプリケーション開発者が期待される出力を最適化するのに役立ちます。"}),e.jsx(n.h3,{children:"パス"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"メッセージID"})},"message_id")}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsxs(n.p,{children:["高評価は",e.jsx(n.code,{children:"like"}),"、低評価は",e.jsx(n.code,{children:"dislike"}),"、高評価の取り消しは",e.jsx(n.code,{children:"null"})]})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"開発者のルールで定義されたユーザー識別子。アプリケーション内で一意である必要があります。"})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"メッセージのフィードバックです。"})},"content")]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"アプリのメッセージの「いいね」とフィードバックを取得",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アプリのエンドユーザーからのフィードバックや「いいね」を取得します。"}),e.jsx(n.h3,{children:"クエリ"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（任意）ページ番号。デフォルト値：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（任意）1ページあたりの件数。デフォルト値：20"})},"limit")}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (リスト) このアプリの「いいね」とフィードバックの一覧を返します。"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`    {
    "data": [
        {
            "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
            "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
            "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
            "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
            "rating": "like",
            "content": "message feedback information-3",
            "from_source": "user",
            "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
            "from_account_id": null,
            "created_at": "2025-04-24T09:24:38",
            "updated_at": "2025-04-24T09:24:38"
        }
    ]
    }
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"テキストから音声",name:"#text-to-audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"テキストを音声に変換します。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"Gofyが生成したテキストメッセージの場合、生成されたmessage-idを直接渡すだけです。バックエンドはmessage-idを使用して対応するコンテンツを検索し、音声情報を直接合成します。message_idとtextの両方が同時に提供された場合、message_idが優先されます。"})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"音声生成コンテンツ。"})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"開発者が定義したユーザー識別子。アプリ内で一意性を確保する必要があります。"})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/text-to-audio",targetCode:`curl -o text-to-audio.mp3 -X POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290",
  "text": "Hello Gofy",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"アプリケーションの基本情報を取得",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"このアプリケーションの基本情報を取得するために使用されます"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) アプリケーションの名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) アプリケーションの説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) アプリケーションのタグ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) アプリケーションのモード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"author_name"})," (string) 作者の名前"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"アプリケーションのパラメータ情報を取得",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ページ開始時に、機能、入力パラメータ名、タイプ、デフォルト値などの情報を取得するために使用されます。"}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) 開始文"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) 開始時の提案質問リスト"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) 回答後の提案質問を有効にします。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) 音声からテキスト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) 引用と帰属",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) 注釈付き返信",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) ユーザー入力フォーム設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) テキスト入力コントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) 段落テキスト入力コントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) ドロップダウンコントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) オプション値"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) ファイルアップロード設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) ドキュメント設定
現在サポートされているドキュメントタイプ：`,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) ドキュメント数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) 画像設定
現在サポートされている画像タイプ：`,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 画像数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) オーディオ設定
現在サポートされているオーディオタイプ：`,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) オーディオ数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) ビデオ設定
現在サポートされているビデオタイプ：`,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) ビデオ数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) カスタム設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) カスタム数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) システムパラメータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) ドキュメントアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) 画像ファイルアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) 音声ファイルアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) 動画ファイルアップロードサイズ制限（MB）"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "opening_statement": "Hello!",
  "suggested_questions_after_answer": {
      "enabled": true
  },
  "speech_to_text": {
      "enabled": true
  },
  "retriever_resource": {
      "enabled": true
  },
  "annotation_reply": {
      "enabled": true
  },
  "user_input_form": [
      {
          "paragraph": {
              "label": "Query",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"アプリのWebApp設定を取得",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アプリの WebApp 設定を取得するために使用します。"}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp 名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) チャットの色テーマ、16 進数形式"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) チャットの色テーマを反転するかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) アイコンタイプ、",e.jsx(n.code,{children:"emoji"}),"-絵文字、",e.jsx(n.code,{children:"image"}),"-画像"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) アイコン。",e.jsx(n.code,{children:"emoji"}),"タイプの場合は絵文字、",e.jsx(n.code,{children:"image"}),"タイプの場合は画像 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) 16 進数形式の背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) アイコンの URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) 著作権情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) プライバシーポリシーのリンク"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) カスタム免責事項"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) デフォルト言語"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) ワークフローの詳細を表示するかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) WebApp のアイコンをチャット内の🤖に置き換えるかどうか"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{})]})}function $n(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(xe,{...i})}):xe(i)}function je(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"文本生成型应用 API"}),`
`,e.jsx(n.p,{children:"文本生成应用无会话支持，适合用于翻译/文章写作/总结 AI 等等。"}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"基础 URL"}),e.jsx(s,{title:"Code",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"鉴权"}),e.jsxs(n.p,{children:["Service API 使用 ",e.jsx(n.code,{children:"API-Key"}),` 进行鉴权。
`,e.jsx("i",{children:e.jsxs(n.strong,{children:["强烈建议开发者把 ",e.jsx(n.code,{children:"API-Key"})," 放在后端存储，而非分享或者放在客户端存储，以免 ",e.jsx(n.code,{children:"API-Key"})," 泄露，导致财产损失。"]})}),`
所有 API 请求都应在 `,e.jsx(n.strong,{children:e.jsx(n.code,{children:"Authorization"})})," HTTP Header 中包含您的 ",e.jsx(n.code,{children:"API-Key"}),"，如下所示："]}),e.jsx(s,{title:"Code",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/completion-messages",method:"POST",title:"发送消息",name:"#Create-Completion-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"发送请求给文本生成型应用。"}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsxs(d,{name:"inputs",type:"object",children:[e.jsx(n.p,{children:`(选填)允许传入 App 定义的各变量值。
inputs 参数包含了多组键值对（Key/Value pairs），每组的键对应一个特定变量，每组的值则是该变量的具体值。
文本生成型应用要求至少传入一组键值对。`}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"}),` (string) 必填
用户输入的文本内容。`]}),`
`]})]},"inputs"),e.jsx(d,{name:"response_mode",type:"string",children:e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," 流式模式（推荐）。基于 SSE（",e.jsx(n.strong,{children:e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"})}),"）实现类似打字机输出方式的流式返回。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` 阻塞模式，等待执行完毕后返回结果。（请求若流程较长可能会被中断）。
`,e.jsx("i",{children:"由于 Cloudflare 限制，请求会在 100 秒超时无返回后中断。"})]}),`
`]})},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`用户标识，用于定义终端用户的身份，方便检索、统计。
由开发者定义规则，需保证用户标识在应用内唯一。`})},"user"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"文件列表，适用于传入文件结合文本理解并回答问题，仅当模型支持 Vision/Video 能力时可用。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 支持类型：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," 具体类型包含：'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," 具体类型包含：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," 具体类型包含：'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," 具体类型包含：'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," 具体类型包含：其他文件类型"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string)  传递方式:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": 文件地址。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": 上传文件。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," 文件地址。（仅当传递方式为 ",e.jsx(n.code,{children:"remote_url"})," 时）。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," 上传文件 ID。（仅当传递方式为 ",e.jsx(n.code,{children:"local_file "}),"时）。"]}),`
`]})]},"files")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(t,{children:[e.jsxs(n.p,{children:["当 ",e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"blocking"}),` 时，返回 ChatCompletionResponse object。
当 `,e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"streaming"}),"时，返回 ChunkChatCompletionResponse object 流式序列。"]}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["返回完整的 App 结果，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"application/json"}),"。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) App 模式，固定为 chat"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 完整回复内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) 模型用量信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用和归属分段列表"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 消息创建时间戳，如：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["返回 App 输出的流式块，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"text/event-stream"}),`。
每个流式块均为 data: 开头，块之间以 `,e.jsx(n.code,{children:"\\n\\n"})," 即两个换行符分隔，如下所示："]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsxs(n.p,{children:["流式块中根据 ",e.jsx(n.code,{children:"event"})," 不同，结构也不同："]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLM 返回文本块事件，即：完整的文本以分块的方式输出。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLM 返回文本块内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," 消息结束事件，收到此事件则代表文本流式返回结束。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) 模型用量信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用和归属分段列表"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS 音频流事件，即：语音合成输出。内容是Mp3格式的音频块，使用 base64 编码后的字符串，播放的时候直接解码即可。(开启自动播放才有此消息)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 语音合成之后的音频块使用 Base64 编码之后的文本内容，播放的时候直接 base64 解码送入播放器即可"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS 音频流结束事件，收到这个事件表示音频流返回结束。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 结束事件是没有音频的，所以这里是空字符串"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` 消息内容替换事件。
开启内容审查和审查输出内容时，若命中了审查条件，则会通过此事件替换消息内容为预设回复。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 替换内容（直接替换 LLM 所有回复文本）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
流式输出过程中出现的异常会以 stream event 形式输出，收到异常事件后即结束。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTP 状态码"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) 错误码"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) 错误消息"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 每 10s 一次的 ping 事件，保持连接存活。"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404，对话不存在"}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"invalid_param"}),"，传入参数异常"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"app_unavailable"}),"，App 配置不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_not_initialize"}),"，无可用模型凭据配置"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_quota_exceeded"}),"，模型调用额度不足"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"model_currently_not_support"}),"，当前模型不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"completion_request_error"}),"，文本生成失败"]}),`
`,e.jsx(n.li,{children:"500，服务内部异常"}),`
`]})]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/completion-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/completion-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": {"query": "Hello, world!"},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(n.h3,{children:"blocking"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "0b089b9a-24d9-48cc-94f8-762677276261",
  "answer": "how are you?",
  "created_at": 1679586667
}
`})})}),e.jsx(n.h3,{children:"streaming"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " I", "created_at": 1679586595}
  data: {"id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "answer": " I", "created_at": 1679586595}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"上传文件",name:"#files-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:[`上传文件（目前仅支持图片）并在发送消息时使用，可实现图文多模态理解。
支持 png, jpg, jpeg, webp, gif 格式。
`,e.jsx("i",{children:"上传的文件仅供当前终端用户使用。"})]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["该接口需使用  ",e.jsx(n.code,{children:"multipart/form-data"})," 进行请求。"]}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsx(n.p,{children:"要上传的文件。"})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，用于定义终端用户的身份，必须和发送消息接口传入 user 保持一致。"})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"成功上传后，服务器会返回文件的 ID 和相关信息。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 文件名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) 文件大小（byte）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) 文件后缀"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) 文件 mime-type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) 上传人 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 上传时间"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"no_file_uploaded"}),"，必须提供文件"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"too_many_files"}),"，目前只接受一个文件"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"unsupported_preview"}),"，该文件不支持预览"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"unsupported_estimate"}),"，该文件不支持估算"]}),`
`,e.jsxs(n.li,{children:["413，",e.jsx(n.code,{children:"file_too_large"}),"，文件太大"]}),`
`,e.jsxs(n.li,{children:["415，",e.jsx(n.code,{children:"unsupported_file_type"}),"，不支持的扩展名，当前只接受文档类文件"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_connection_failed"}),"，无法连接到 S3 服务"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_permission_denied"}),"，无权限上传文件到 S3"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_file_too_large"}),"，文件超出 S3 大小限制"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": 123,
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"获取终端用户",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"通过终端用户 ID 获取终端用户信息。"}),e.jsxs(n.p,{children:["当其他 API 返回终端用户 ID（例如：上传文件接口返回的 ",e.jsx(n.code,{children:"created_by"}),"）时，可使用该接口查询对应的终端用户信息。"]}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) 必需
终端用户 ID。`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"返回 EndUser 对象。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) 工作空间（Tenant）ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) 应用 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 终端用户类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) 外部用户 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) 是否匿名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 时间"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404，",e.jsx(n.code,{children:"end_user_not_found"}),"，终端用户不存在"]}),`
`,e.jsx(n.li,{children:"500，内部服务器错误"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"文件预览",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"预览或下载已上传的文件。此端点允许您访问先前通过文件上传 API 上传的文件。"}),e.jsx("i",{children:"文件只能在属于请求应用程序的消息范围内访问。"}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) 必需
要预览的文件的唯一标识符，从文件上传 API 响应中获得。`]}),`
`]}),e.jsx(n.h3,{children:"查询参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) 可选
是否强制将文件作为附件下载。默认为 `,e.jsx(n.code,{children:"false"}),"（在浏览器中预览）。"]}),`
`]}),e.jsx(n.h3,{children:"响应"}),e.jsx(n.p,{children:"返回带有适当浏览器显示或下载标头的文件内容。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," 根据文件 MIME 类型设置"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," 文件大小（以字节为单位，如果可用）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"})," 如果 ",e.jsx(n.code,{children:"as_attachment=true"}),' 则设置为 "attachment"']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," 用于性能的缓存标头"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' 对于音频/视频文件设置为 "bytes"']}),`
`]}),e.jsx(n.h3,{children:"错误"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", 参数输入异常"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", 文件访问被拒绝或文件不属于当前应用程序"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", 文件未找到或已被删除"]}),`
`,e.jsx(n.li,{children:"500, 服务内部错误"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"请求示例"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"作为附件下载"}),e.jsx(s,{title:"下载请求",tag:"GET",label:"/files/:file_id/preview?as_attachment=true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"响应标头示例"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"文件下载响应标头"}),e.jsx(s,{title:"Download Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/completion-messages/:task_id/stop",method:"POST",title:"停止响应",name:"#Stop"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"仅支持流式模式。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，可在流式返回 Chunk 中获取"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
用户标识，用于定义终端用户的身份，必须和发送消息接口传入 user 保持一致。API 无法访问 WebApp 创建的会话。`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"})," (string) 固定返回 success"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/completion-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/completion-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{ "user": "abc-123"}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"消息反馈（点赞）",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"消息终端用户反馈、点赞，方便应用开发者优化输出预期。"}),e.jsx(n.h3,{children:"Path Params"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"消息 ID"})},"message_id")}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsx(n.p,{children:"点赞 like, 点踩 dislike,  撤销点赞 null"})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"消息反馈的具体信息。"})},"content")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"})," (string) 固定返回 success"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"Get feedbacks of application",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Get application's feedbacks."}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（optional）pagination，default：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（optional） records per page default：20"})},"limit")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (List) return apps feedback list."]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`  {
      "data": [
          {
              "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
              "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
              "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
              "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
              "rating": "like",
              "content": "message feedback information-3",
              "from_source": "user",
              "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
              "from_account_id": null,
              "created_at": "2025-04-24T09:24:38",
              "updated_at": "2025-04-24T09:24:38"
          }
      ]
  }
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"文字转语音",name:"#text-to-audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"文字转语音。"}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"Gofy 生成的文本消息，那么直接传递生成的message-id 即可，后台会通过 message_id 查找相应的内容直接合成语音信息。如果同时传 message_id 和 text，优先使用 message_id。"})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"语音生成内容。如果没有传 message-id的话，则会使用这个字段的内容"})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/text-to-audio",targetCode:`curl -o text-to-audio.mp3 -X POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290",
  "text": "你好Gofy",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"获取应用基本信息",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取应用的基本信息"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 应用名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 应用描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) 应用标签"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) 应用模式"]}),`
`,e.jsx(n.li,{children:"'author_name' (string) 作者名称"}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"获取应用参数",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于进入页面一开始，获取功能开关、输入参数名称、类型及默认值等使用。"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) 开场白"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) 开场推荐问题列表"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) 启用回答后给出推荐问题。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) 语音转文本",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) 引用和归属",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) 标记回复",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) 用户输入表单配置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) 文本输入控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) 段落文本输入控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) 下拉控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) 选项值"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) 文件上传配置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) 文档设置
当前仅支持文档类型：`,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 文档数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) 图片设置
当前仅支持图片类型：`,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 图片数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) 音频设置
当前仅支持音频类型：`,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 音频数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) 视频设置
当前仅支持视频类型：`,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 视频数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) 自定义设置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 自定义数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) 系统参数",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) 文档上传大小限制 (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) 图片文件上传大小限制（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) 音频文件上传大小限制 (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) 视频文件上传大小限制 (MB)"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'\\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "introduction": "nice to meet you",
  "user_input_form": [
    {
      "text-input": {
        "label": "a",
        "variable": "a",
        "required": true,
        "max_length": 48,
        "default": ""
      }
    },
    {
      // ...
    }
  ],
  "file_upload": {
    "image": {
      "enabled": true,
      "number_limits": 3,
      "transfer_methods": [
        "remote_url",
        "local_file"
      ]
    }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"获取应用 WebApp 设置",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取应用的 WebApp 设置"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp 名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) 聊天颜色主题，hex 格式"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) 聊天颜色主题是否反转"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) 图标类型，",e.jsx(n.code,{children:"emoji"}),"-表情，",e.jsx(n.code,{children:"image"}),"-图片"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) 图标，如果是 ",e.jsx(n.code,{children:"emoji"})," 类型，则是 emoji 表情符号，如果是 ",e.jsx(n.code,{children:"image"})," 类型，则是图片 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) hex 格式的背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) 图标 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) 版权信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) 隐私政策链接"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) 自定义免责声明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) 默认语言"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) 是否显示工作流详情"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) 是否使用 WebApp 图标替换聊天中的 🤖"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations",method:"GET",title:"获取标注列表",name:"#annotation_list"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"页码"})},"page"),e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"每页数量"})},"limit")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/apps/annotations",targetCode:`curl --location --request GET '${i.apiBaseUrl}/apps/annotations?page=1&limit=20' \\
--header 'Authorization: Bearer {api_key}'`,children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl --location --request GET '\${props.apiBaseUrl}/apps/annotations?page=1&limit=20' \\
--header 'Authorization: Bearer {api_key}'
`})})}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "data": [
    {
      "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
      "question": "What is your name?",
      "answer": "I am Gofy.",
      "hit_count": 0,
      "created_at": 1735625869
    }
  ],
  "has_more": false,
  "limit": 20,
  "total": 1,
  "page": 1
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations",method:"POST",title:"创建标注",name:"#create_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"question",type:"string",children:e.jsx(n.p,{children:"问题"})},"question"),e.jsx(d,{name:"answer",type:"string",children:e.jsx(n.p,{children:"答案内容"})},"answer")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/apps/annotations",targetCode:`curl --location --request POST '${i.apiBaseUrl}/apps/annotations' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{"question": "What is your name?","answer": "I am Gofy."}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
  "question": "What is your name?",
  "answer": "I am Gofy.",
  "hit_count": 0,
  "created_at": 1735625869
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations/{annotation_id}",method:"PUT",title:"更新标注",name:"#update_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"annotation_id",type:"string",children:e.jsx(n.p,{children:"标注 ID"})},"annotation_id"),e.jsx(d,{name:"question",type:"string",children:e.jsx(n.p,{children:"问题"})},"question"),e.jsx(d,{name:"answer",type:"string",children:e.jsx(n.p,{children:"答案内容"})},"answer")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/apps/annotations/{annotation_id}",targetCode:`curl --location --request PUT '${i.apiBaseUrl}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{"question": "What is your name?","answer": "I am Gofy."}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
  "question": "What is your name?",
  "answer": "I am Gofy.",
  "hit_count": 0,
  "created_at": 1735625869
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations/{annotation_id}",method:"DELETE",title:"删除标注",name:"#delete_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"annotation_id",type:"string",children:e.jsx(n.p,{children:"标注 ID"})},"annotation_id")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/apps/annotations/{annotation_id}",targetCode:`curl --location --request DELETE '${i.apiBaseUrl}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json'`,children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl --location --request DELETE '\${props.apiBaseUrl}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}'
`})})}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotation-reply/{action}",method:"POST",title:"标注回复初始设置",name:"#initial_annotation_reply_settings"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"action",type:"string",children:e.jsx(n.p,{children:"动作，只能是 'enable' 或 'disable'"})},"action"),e.jsx(d,{name:"embedding_provider_name",type:"string",children:e.jsx(n.p,{children:"指定的嵌入模型提供商，必须先在系统内设定好接入的模型，对应的是 provider 字段"})},"embedding_provider_name"),e.jsx(d,{name:"embedding_model_name",type:"string",children:e.jsx(n.p,{children:"指定的嵌入模型，对应的是 model 字段"})},"embedding_model_name"),e.jsx(d,{name:"score_threshold",type:"number",children:e.jsx(n.p,{children:"相似度阈值，当相似度大于该阈值时，系统会自动回复，否则不回复"})},"score_threshold")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.p,{children:`嵌入模型的提供商和模型名称可以通过以下接口获取：v1/workspaces/current/models/model-types/text-embedding，具体见：通过 API 维护知识库。使用的 Authorization 是 Dataset 的 API Token。
该接口是异步执行，所以会返回一个 job_id，通过查询 job 状态接口可以获取到最终的执行结果。`}),e.jsx(s,{title:"Request",tag:"POST",label:"/apps/annotation-reply/{action}",targetCode:`curl --location --request POST '${i.apiBaseUrl}/apps/annotation-reply/{action}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{"score_threshold": 0.9, "embedding_provider_name": "zhipu", "embedding_model_name": "embedding_3"}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "job_id": "b15c8f68-1cf4-4877-bf21-ed7cf2011802",
  "job_status": "waiting"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotation-reply/{action}/status/{job_id}",method:"GET",title:"查询标注回复初始设置任务状态",name:"#initial_annotation_reply_settings_task_status"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"action",type:"string",children:e.jsx(n.p,{children:"动作，只能是 'enable' 或 'disable'，并且必须和标注回复初始设置接口的动作一致"})},"action"),e.jsx(d,{name:"job_id",type:"string",children:e.jsx(n.p,{children:"任务 ID，从标注回复初始设置接口返回的 job_id"})},"job_id")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/apps/annotations",targetCode:`curl --location --request GET '${i.apiBaseUrl}/apps/annotation-reply/{action}/status/{job_id}' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "job_id": "b15c8f68-1cf4-4877-bf21-ed7cf2011802",
  "job_status": "waiting",
  "error_msg": ""
}
`})})})]})]})]})}function Hn(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(je,{...i})}):je(i)}function ue(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"Advanced Chat App API"}),`
`,e.jsx(n.p,{children:"Chat applications support session persistence, allowing previous chat history to be used as context for responses. This can be applicable for chatbot, customer service AI, etc."}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"Base URL"}),e.jsx(s,{title:"Code",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"Authentication"}),e.jsxs(n.p,{children:["The Service API uses ",e.jsx(n.code,{children:"API-Key"}),` authentication.
`,e.jsx("i",{children:e.jsx(n.strong,{children:"Strongly recommend storing your API Key on the server-side, not shared or stored on the client-side, to avoid possible API-Key leakage that can lead to serious consequences."})})]}),e.jsxs(n.p,{children:["For all API requests, include your API Key in the ",e.jsx(n.code,{children:"Authorization"}),"HTTP Header, as shown below:"]}),e.jsx(s,{title:"Code",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages",method:"POST",title:"Send Chat Message",name:"#Send-Chat-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Send a request to the chat application."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"query",type:"string",children:e.jsx(n.p,{children:"User Input/Question content"})},"query"),e.jsx(d,{name:"inputs",type:"object",children:e.jsxs(n.p,{children:[`Allows the entry of various variable values defined by the App.
The `,e.jsx(n.code,{children:"inputs"}),` parameter contains multiple key/value pairs, with each key corresponding to a specific variable and each value being the specific value for that variable.
If the variable is of file type, specify an object that has the keys described in `,e.jsx(n.code,{children:"files"}),` below.
Default `,e.jsx(n.code,{children:"{}"})]})},"inputs"),e.jsxs(d,{name:"response_mode",type:"string",children:[e.jsx(n.p,{children:"The mode of response return, supporting:"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," Streaming mode (recommended), implements a typewriter-like output through SSE (",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` Blocking mode, returns result after execution is complete. (Requests may be interrupted if the process is long)
Due to Cloudflare restrictions, the request will be interrupted without a return after 100 seconds.`]}),`
`]})]},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application. The Service API does not share conversations created by the WebApp.`})},"user"),e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"Conversation ID, to continue the conversation based on previous chat records, it is necessary to pass the previous message's conversation_id."})},"conversation_id"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"File list, suitable for inputting files combined with text understanding and answering questions, available only when the model supports Vision/Video capability."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) Supported type:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," Supported types include: 'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," Supported types include: 'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," Supported types include: 'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," Supported types include: 'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," Supported types include: other file types"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) Transfer method:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": File URL."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": Upload file."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," File URL. (Only when transfer method is ",e.jsx(n.code,{children:"remote_url"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," Upload file ID. (Only when transfer method is ",e.jsx(n.code,{children:"local_file"}),")."]}),`
`]})]},"files"),e.jsx(d,{name:"auto_generate_name",type:"bool",children:e.jsxs(n.p,{children:["Auto-generate title, default is ",e.jsx(n.code,{children:"true"}),`.
If set to `,e.jsx(n.code,{children:"false"}),", can achieve async title generation by calling the conversation rename API and setting ",e.jsx(n.code,{children:"auto_generate"})," to ",e.jsx(n.code,{children:"true"}),"."]})},"auto_generate_name"),e.jsx(d,{name:"workflow_id",type:"string",children:e.jsx(n.p,{children:"(Optional) Workflow ID to specify a specific version, if not provided, uses the default published version."})},"workflow_id"),e.jsxs(d,{name:"trace_id",type:"string",children:[e.jsxs(n.p,{children:["(Optional) Trace ID. Used for integration with existing business trace components to achieve end-to-end distributed tracing. If not provided, the system will automatically generate a trace_id. Supports the following three ways to pass, in order of priority:",e.jsx("br",{})]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["Header: via HTTP Header ",e.jsx("code",{children:"X-Trace-Id"}),", highest priority.",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["Query parameter: via URL query parameter ",e.jsx("code",{children:"trace_id"}),".",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["Request Body: via request body field ",e.jsx("code",{children:"trace_id"})," (i.e., this field).",e.jsx("br",{})]}),`
`]})]},"trace_id")]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:`When response_mode is blocking, return a CompletionResponse object.
When response_mode is streaming, return a ChunkCompletionResponse stream.`}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["Returns the complete App result, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"application/json"}),"."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) Event type, fixed to ",e.jsx(n.code,{children:"message"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) unique ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) App mode, fixed as ",e.jsx(n.code,{children:"chat"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) Complete response content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) Metadata",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) Model usage information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) Citation and Attribution List"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Message creation timestamp, e.g., 1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["Returns the stream chunks outputted by the App, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"text/event-stream"}),`.
Each streaming chunk starts with `,e.jsx(n.code,{children:"data:"}),", separated by two newline characters ",e.jsx(n.code,{children:"\\n\\n"}),", as shown below:"]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsxs(n.p,{children:["The structure of the streaming chunks varies depending on the ",e.jsx(n.code,{children:"event"}),":"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLM returns text chunk event, i.e., the complete text is output in a chunked fashion.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLM returned text chunk content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_file"})," Message file event, a new file has created by tool",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) File unique ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"}),' (string) File type，only allow "image" currently']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) Belongs to, it will only be an 'assistant' here"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) Remote url of file"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"}),"  (string) Conversation ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," Message end event, receiving this event means streaming has ended.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) Metadata",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) Model usage information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) Citation and Attribution List"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS audio stream event, that is, speech synthesis output. The content is an audio block in Mp3 format, encoded as a base64 string. When playing, simply decode the base64 and feed it into the player. (This message is available only when auto-play is enabled)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the stop response interface below"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The audio after speech synthesis, encoded in base64 text content, when playing, simply decode the base64 and feed it into the player"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g.: 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS audio stream end event, receiving this event indicates the end of the audio stream.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the stop response interface below"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The end event has no audio, so this is an empty string"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g.: 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` Message content replacement event.
When output content moderation is enabled, if the content is flagged, then the message content will be replaced with a preset reply through this event.`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) Replacement content (directly replaces all LLM reply text)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," workflow starts execution",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"workflow_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," node execution started",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"node_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ID of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) type of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) name of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) Execution sequence number, used to display Tracing Node sequence"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional Prefix node ID, used for canvas display execution path"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) Contents of all preceding node variables used in the node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) timestamp of start, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," node execution ends, success or failure in different states in the same event",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"node_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ID of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) type of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) name of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) Execution sequence number, used to display Tracing Node sequence"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional Prefix node ID, used for canvas display execution path"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) Contents of all preceding node variables used in the node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) Optional node process data"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) meta data",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) optional tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) optional Total cost"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) optional e.g. ",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) timestamp of start, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," workflow execution ends, success or failure in different states in the same event",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"workflow_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) default 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) start time"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) end time"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
Exceptions that occur during the streaming process will be output in the form of stream events, and reception of an error event will end the stream.`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTP status code"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) Error code"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) Error message"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," Ping event every 10 seconds to keep the connection alive."]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404, Conversation does not exists"}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", abnormal parameter input"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", App configuration unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", no available model credential configuration"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", model invocation quota insufficient"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", current model unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_not_found"}),", specified workflow version not found"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"draft_workflow_error"}),", cannot use draft workflow version"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_id_format_error"}),", invalid workflow_id format, expected UUID format"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"completion_request_error"}),", text generation failed"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/chat-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "query": "What are the specs of the iPhone 13 Pro Max?",
  "response_mode": "streaming",
  "conversation_id": "",
  "user": "abc-123",
  "files": [
      {
          "type": "image",
          "transfer_method": "remote_url",
          "url": "https://cloud.gofy.ai/logo/logo-site.png"
      }
  ]
}'`}),e.jsx(n.h3,{children:"Blocking Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "event": "message",
    "task_id": "c3800678-a077-43df-a102-53f23ed20b88",
    "id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "message_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
    "mode": "chat",
    "answer": "iPhone 13 Pro Max specs are listed here:...",
    "metadata": {
        "usage": {
            "prompt_tokens": 1033,
            "prompt_unit_price": "0.001",
            "prompt_price_unit": "0.001",
            "prompt_price": "0.0010330",
            "completion_tokens": 128,
            "completion_unit_price": "0.002",
            "completion_price_unit": "0.001",
            "completion_price": "0.0002560",
            "total_tokens": 1161,
            "total_price": "0.0012890",
            "currency": "USD",
            "latency": 0.7682376249867957
        },
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ]
    },
    "created_at": 1705407629
}
`})})}),e.jsx(n.h3,{children:"Streaming Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " I", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": "'m", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " glad", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " to", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " meet", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " you", "created_at": 1679586595}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"File Upload",name:"#file-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:`Upload a file for use when sending messages, enabling multimodal understanding of images and text.
Supports any formats that are supported by your application.
Uploaded files are for use by the current end-user only.`}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["This interface requires a ",e.jsx(n.code,{children:"multipart/form-data"})," request."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file"}),` (File) Required
The file to be uploaded.`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, defined by the developer's rules, must be unique within the application. The Service API does not share conversations created by the WebApp.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"After a successful upload, the server will return the file's ID and related information."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) File name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) File size (bytes)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) File extension"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) File mime-type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) End-user ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"no_file_uploaded"}),", a file must be provided"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"too_many_files"}),", currently only one file is accepted"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_preview"}),", the file does not support preview"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_estimate"}),", the file does not support estimation"]}),`
`,e.jsxs(n.li,{children:["413, ",e.jsx(n.code,{children:"file_too_large"}),", the file is too large"]}),`
`,e.jsxs(n.li,{children:["415, ",e.jsx(n.code,{children:"unsupported_file_type"}),", unsupported extension, currently only document files are accepted"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_connection_failed"}),", unable to connect to S3 service"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_permission_denied"}),", no permission to upload files to S3"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_file_too_large"}),", file exceeds S3 size limit"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"Get End User",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve an end user by ID."}),e.jsxs(n.p,{children:["This is useful when other APIs return an end-user ID (e.g. ",e.jsx(n.code,{children:"created_by"})," from File Upload)."]}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) Required
End user ID.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns an EndUser object."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) Tenant ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) App ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) End user type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) External user ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) Whether anonymous"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) Session ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 datetime"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 datetime"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"end_user_not_found"}),", end user not found"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"File Preview",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Preview or download uploaded files. This endpoint allows you to access files that have been previously uploaded via the File Upload API."}),e.jsx("i",{children:"Files can only be accessed if they belong to messages within the requesting application."}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) Required
The unique identifier of the file to preview, obtained from the File Upload API response.`]}),`
`]}),e.jsx(n.h3,{children:"Query Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) Optional
Whether to force download the file as an attachment. Default is `,e.jsx(n.code,{children:"false"})," (preview in browser)."]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns the file content with appropriate headers for browser display or download."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," Set based on file mime type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," File size in bytes (if available)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"}),' Set to "attachment" if ',e.jsx(n.code,{children:"as_attachment=true"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," Caching headers for performance"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' Set to "bytes" for audio/video files']}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", abnormal parameter input"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", file access denied or file does not belong to current application"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", file not found or has been deleted"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Download as Attachment"}),e.jsx(s,{title:"Download Request",tag:"GET",label:"/files/:file_id/preview?as_attachment=true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"Response Headers Example"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"Download Response Headers"}),e.jsx(s,{title:"Download Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages/:task_id/stop",method:"POST",title:"Stop Generate",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Only supported in streaming mode."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, can be obtained from the streaming chunk return"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, used to define the identity of the end-user, must be consistent with the user passed in the message sending interface. The Service API does not share conversations created by the WebApp.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/chat-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{
  "user": "abc-123"
}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"Message Feedback",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"End-users can provide feedback messages, facilitating application developers to optimize expected outputs."}),e.jsx(n.h3,{children:"Path"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"Message ID"})},"message_id")}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsxs(n.p,{children:["Upvote as ",e.jsx(n.code,{children:"like"}),", downvote as ",e.jsx(n.code,{children:"dislike"}),", revoke upvote as ",e.jsx(n.code,{children:"null"})]})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"User identifier, defined by the developer's rules, must be unique within the application. The Service API does not share conversations created by the WebApp."})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"The specific content of message feedback."})},"content")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"Get feedbacks of application",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Get application's feedbacks."}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（optional）pagination，default：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（optional） records per page default：20"})},"limit")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (List) return apps feedback list."]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`  {
      "data": [
          {
              "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
              "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
              "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
              "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
              "rating": "like",
              "content": "message feedback information-3",
              "from_source": "user",
              "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
              "from_account_id": null,
              "created_at": "2025-04-24T09:24:38",
              "updated_at": "2025-04-24T09:24:38"
          }
      ]
  }
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/{message_id}/suggested",method:"GET",title:"Next Suggested Questions",name:"#suggested"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Get next questions suggestions for the current message"}),e.jsx(n.h3,{children:"Path Params"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"Message ID"})},"message_id")}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`})},"user")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/messages/{message_id}/suggested",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/messages/{message_id}/suggested?user=abc-123' \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success",
  "data": [
        "a",
        "b",
        "c"
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages",method:"GET",title:"Get Conversation History Messages",name:"#messages"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:["Returns historical chat records in a scrolling load format, with the first page returning the latest ",e.jsx(n.code,{children:"{limit}"})," messages, i.e., in reverse order."]}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"Conversation ID"})},"conversation_id"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`})},"user"),e.jsx(d,{name:"first_id",type:"string",children:e.jsx(n.p,{children:"The ID of the first chat record on the current page, default is null."})},"first_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"How many chat history messages to return in one request, default is 20."})},"limit")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) Message list",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) User input parameters."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"})," (string) User input / question content."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[object]) Message files",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) File type, image for images"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) File preview URL, use the File Preview API (",e.jsx(n.code,{children:"/files/{file_id}/preview"}),") to access the file"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) belongs to，user orassistant"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) Response message content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"feedback"})," (object) Feedback information",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"rating"})," (string) Upvote as ",e.jsx(n.code,{children:"like"})," / Downvote as ",e.jsx(n.code,{children:"dislike"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) Citation and Attribution List"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) Whether there is a next page"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) Number of returned items, if input exceeds system limit, returns system limit amount"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/messages",targetCode:`curl -X GET '${i.appDetail.api_base_url}/messages?user=abc-123&conversation_id={conversation_id}'
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
        "id": "a076a87f-31e5-48dc-b452-0061adbbc922",
        "conversation_id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
        "inputs": {
            "name": "gofy"
        },
        "query": "iphone 13 pro",
        "answer": "The iPhone 13 Pro, released on September 24, 2021, features a 6.1-inch display with a resolution of 1170 x 2532. It is equipped with a Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard) processor, 6 GB of RAM, and offers storage options of 128 GB, 256 GB, 512 GB, and 1 TB. The camera is 12 MP, the battery capacity is 3095 mAh, and it runs on iOS 15.",
        "message_files": [],
        "feedback": null,
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ],
        "created_at": 1705569239,
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations",method:"GET",title:"Get Conversations",name:"#conversations"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve the conversation list for the current user, defaulting to the most recent 20 entries."}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"(Optional) The ID of the last record on the current page, default is null."})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"(Optional) How many records to return in one request, default is the most recent 20 entries. Maximum 100, minimum 1."})},"limit"),e.jsxs(d,{name:"sort_by",type:"string",children:[e.jsx(n.p,{children:"(Optional) Sorting Field, Default: -updated_at (sorted in descending order by update time)"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"Available Values: created_at, -created_at, updated_at, -updated_at"}),`
`,e.jsx(n.li,{children:'The symbol before the field represents the order or reverse, "-" represents reverse order.'}),`
`]})]},"sort_by")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) List of conversations",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Conversation name, by default, is generated by LLM."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) User input parameters."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) Conversation status"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) Introduction"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) Update timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) Number of entries returned, if input exceeds system limit, system limit number is returned"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/conversations",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations?user=abc-123&last_id=&limit=20' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
      "id": "10799fb8-64f7-4296-bbf7-b42bfbe0ae54",
      "name": "New chat",
      "inputs": {
          "book": "book",
          "myName": "Lucy"
      },
      "status": "normal",
      "created_at": 1679667915,
      "updated_at": 1679667915
    },
    {
      "id": "hSIhXBhNe8X1d8Et"
      // ...
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id",method:"DELETE",title:"Delete Conversation",name:"#delete"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Delete a conversation."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the application."})},"user")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"DELETE",label:"/conversations/:conversation_id",targetCode:`curl -X DELETE '${i.appDetail.api_base_url}/conversations/{conversation_id}' \\
--header 'Content-Type: application/json' \\
--header 'Accept: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data '{
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/name",method:"POST",title:"Conversation Rename",name:"#rename"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Request Body"}),e.jsx(n.p,{children:"Rename the session, the session name is used for display on clients that support multiple sessions."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`]}),e.jsxs(t,{children:[e.jsx(d,{name:"name",type:"string",children:e.jsxs(n.p,{children:["(Optional) The name of the conversation. This parameter can be omitted if ",e.jsx(n.code,{children:"auto_generate"})," is set to ",e.jsx(n.code,{children:"true"}),"."]})},"name"),e.jsx(d,{name:"auto_generate",type:"bool",children:e.jsxs(n.p,{children:["(Optional) Automatically generate the title, default is ",e.jsx(n.code,{children:"false"})]})},"auto_generate"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the application."})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Conversation name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) User input parameters"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) Conversation status"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) Introduction"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) Update timestamp, e.g., 1705395332"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/conversations/:conversation_id/name",targetCode:`curl -X POST '${i.appDetail.api_base_url}/conversations/{conversation_id}/name' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "name": "",
  "auto_generate": true,
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
    "name": "Chat vs AI",
    "inputs": {},
    "status": "normal",
    "introduction": "",
    "created_at": 1705569238,
    "updated_at": 1705569238
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables",method:"GET",title:"Get Conversation Variables",name:"#conversation-variables"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve variables from a specific conversation. This endpoint is useful for extracting structured data that was captured during the conversation."}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsx(t,{children:e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"The ID of the conversation to retrieve variables from."})},"conversation_id")}),e.jsx(n.h3,{children:"Query Parameters"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the application"})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"(Optional) The ID of the last record on the current page, default is null."})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"(Optional) How many records to return in one request, default is the most recent 20 entries. Maximum 100, minimum 1."})},"limit")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) Number of items per page"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) Whether there is a next page"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) List of variables",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Variable name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) Variable type (string, number, object, etc.)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (string) Variable value"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) Variable description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) Last update timestamp"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", Conversation not found"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/conversations/:conversation_id/variables",debug:"true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Request with variable name filter",language:"bash",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123&variable_name=customer_name' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 100,
  "has_more": false,
  "data": [
    {
      "id": "variable-uuid-1",
      "name": "customer_name",
      "value_type": "string",
      "value": "John Doe",
      "description": "Customer name extracted from the conversation",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    },
    {
      "id": "variable-uuid-2",
      "name": "order_details",
      "value_type": "json",
      "value": "{\\"product\\":\\"Widget\\",\\"quantity\\":5,\\"price\\":19.99}",
      "description": "Order details from the customer",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables/:variable_id",method:"PUT",title:"Update Conversation Variable",name:"#update-conversation-variable"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Update the value of a specific conversation variable. This endpoint allows you to modify the value of a variable that was captured during the conversation while preserving its name, type, and description."}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"The ID of the conversation containing the variable to update."})},"conversation_id"),e.jsx(d,{name:"variable_id",type:"string",children:e.jsx(n.p,{children:"The ID of the variable to update."})},"variable_id")]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"value",type:"any",children:e.jsx(n.p,{children:"The new value for the variable. Must match the variable's expected type (string, number, object, etc.)."})},"value"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the application."})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns the updated variable object with:"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Variable name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) Variable type (string, number, object, etc.)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (any) Updated variable value"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) Variable description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) Last update timestamp"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"Type mismatch: variable expects {expected_type}, but got {actual_type} type"}),", Value type doesn't match variable's expected type"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", Conversation not found"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_variable_not_exists"}),", Variable not found"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/conversations/:conversation_id/variables/:variable_id",targetCode:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": "Updated Value",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Update with different value types",targetCode:[{title:"String",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": "New string value",
  "user": "abc-123"
}'`},{title:"Number",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": 42,
  "user": "abc-123"
}'`},{title:"Object",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": {"product": "Widget", "quantity": 10, "price": 29.99},
  "user": "abc-123"
}'`}]}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "variable-uuid-1",
  "name": "customer_name",
  "value_type": "string",
  "value": "Updated Value",
  "description": "Customer name extracted from the conversation",
  "created_at": 1650000000000,
  "updated_at": 1650000001000
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/audio-to-text",method:"POST",title:"Speech to Text",name:"#audio-to-text"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"This endpoint requires a multipart/form-data request."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsxs(n.p,{children:[`Audio file.
Supported formats: `,e.jsx(n.code,{children:"['mp3', 'mp4', 'mpeg', 'mpga', 'm4a', 'wav', 'webm']"}),`
File size limit: 15MB`]})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"User identifier, defined by the developer's rules, must be unique within the application."})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) Output text"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/audio-to-text",targetCode:`curl -X POST '${i.appDetail.api_base_url}/audio-to-text' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=audio/[mp3|mp4|mpeg|mpga|m4a|wav|webm]'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "text": ""
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"Text to Audio",name:"#text-to-audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Text to speech."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"For text messages generated by Gofy, simply pass the generated message-id directly. The backend will use the message-id to look up the corresponding content and synthesize the voice information directly. If both message_id and text are provided simultaneously, the message_id is given priority."})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"Speech generated content。"})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the app."})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/text-to-audio",targetCode:`curl -o text-to-audio.mp3 -X POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290",
  "text": "Hello Gofy",
  "user": "abc-123",
}'`}),e.jsx(s,{title:"headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"Get Application Basic Information",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get basic information about this application"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) application name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) application description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) application tags"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) application mode"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"author_name"})," (string) application author name"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "advanced-chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"Get Application Parameters Information",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used at the start of entering the page to obtain information such as features, input parameter names, types, and default values."}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) Opening statement"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) List of suggested questions for the opening"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) Suggest questions after enabling the answer.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) Speech to text",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text_to_speech"})," (object) Text to speech",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"voice"})," (string) Voice type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"language"})," (string) Language"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"autoPlay"})," (string) Auto play",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"}),"   Enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"disabled"}),"  Disabled"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) Citation and Attribution",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) Annotation reply",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) User input form configuration",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) Text input control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) Paragraph text input control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) Dropdown control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) Option values"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) File upload configuration",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) Document settings
Currently only supports document types: `,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Document number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) Image settings
Currently only supports image types: `,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Image number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) Audio settings
Currently only supports audio types: `,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Audio number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) Video settings
Currently only supports video types: `,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Video number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) Custom settings",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Custom number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) System parameters",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) Document upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) Image file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) Audio file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) Video file upload size limit (MB)"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:`curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "opening_statement": "Hello!",
  "suggested_questions_after_answer": {
      "enabled": true
  },
  "speech_to_text": {
      "enabled": true
  },
  "text_to_speech": {
      "enabled": true,
      "voice": "sambert-zhinan-v1",
      "language": "zh-Hans",
      "autoPlay": "disabled"
  },
  "retriever_resource": {
      "enabled": true
  },
  "annotation_reply": {
      "enabled": true
  },
  "user_input_form": [
      {
          "paragraph": {
              "label": "Query",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/meta",method:"GET",title:"Get Application Meta Information",name:"#meta"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get icons of tools in this application"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_icons"}),"(object[string]) tool icons",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_name"})," (string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (object|string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["(object) icon object",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"background"})," (string) background color in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"content"}),"(string) emoji"]}),`
`]}),`
`]}),`
`,e.jsx(n.li,{children:"(string) url of icon"}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/meta",targetCode:`curl -X GET '${i.appDetail.api_base_url}/meta' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "tool_icons": {
    "dalle2": "https://cloud.gofy.ai/console/api/workspaces/current/tool-provider/builtin/dalle/icon",
    "api_tool": {
      "background": "#252525",
      "content": "\\ud83d\\ude01"
    }
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"Get Application WebApp Settings",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get the WebApp settings of the application."}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) Chat color theme, in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) Whether the chat color theme is inverted"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) Icon type, ",e.jsx(n.code,{children:"emoji"})," - emoji, ",e.jsx(n.code,{children:"image"})," - picture"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) Icon. If it's ",e.jsx(n.code,{children:"emoji"})," type, it's an emoji symbol; if it's ",e.jsx(n.code,{children:"image"})," type, it's an image URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) Background color in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) Icon URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) Description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) Copyright information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) Privacy policy link"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) Custom disclaimer"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) Default language"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) Whether to show workflow details"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) Whether to replace 🤖 in chat with the WebApp icon"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations",method:"GET",title:"Get Annotation List",name:"#annotation_list"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"Page number"})},"page"),e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"Number of items returned, default 20, range 1-100"})},"limit")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/apps/annotations",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/apps/annotations?page=1&limit=20' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "data": [
    {
      "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
      "question": "What is your name?",
      "answer": "I am Gofy.",
      "hit_count": 0,
      "created_at": 1735625869
    }
  ],
  "has_more": false,
  "limit": 20,
  "total": 1,
  "page": 1
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations",method:"POST",title:"Create Annotation",name:"#create_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"question",type:"string",children:e.jsx(n.p,{children:"Question"})},"question"),e.jsx(d,{name:"answer",type:"string",children:e.jsx(n.p,{children:"Answer"})},"answer")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/apps/annotations",targetCode:`curl --location --request POST '${i.appDetail.api_base_url}/apps/annotations' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "question": "What is your name?",
  "answer": "I am Gofy."
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
  "question": "What is your name?",
  "answer": "I am Gofy.",
  "hit_count": 0,
  "created_at": 1735625869
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations/{annotation_id}",method:"PUT",title:"Update Annotation",name:"#update_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"annotation_id",type:"string",children:e.jsx(n.p,{children:"Annotation ID"})},"annotation_id"),e.jsx(d,{name:"question",type:"string",children:e.jsx(n.p,{children:"Question"})},"question"),e.jsx(d,{name:"answer",type:"string",children:e.jsx(n.p,{children:"Answer"})},"answer")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/apps/annotations/{annotation_id}",targetCode:`curl --location --request PUT '${i.appDetail.api_base_url}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "question": "What is your name?",
  "answer": "I am Gofy."
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
  "question": "What is your name?",
  "answer": "I am Gofy.",
  "hit_count": 0,
  "created_at": 1735625869
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations/{annotation_id}",method:"DELETE",title:"Delete Annotation",name:"#delete_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"annotation_id",type:"string",children:e.jsx(n.p,{children:"Annotation ID"})},"annotation_id")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"DELETE",label:"/apps/annotations/{annotation_id}",targetCode:`curl --location --request DELETE '${i.appDetail.api_base_url}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotation-reply/{action}",method:"POST",title:"Initial Annotation Reply Settings",name:"#initial_annotation_reply_settings"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"action",type:"string",children:e.jsx(n.p,{children:"Action, can only be 'enable' or 'disable'"})},"action"),e.jsx(d,{name:"embedding_provider_name",type:"string",children:e.jsx(n.p,{children:"Specified embedding model provider, must be set up in the system first, corresponding to the provider field(Optional)"})},"embedding_provider_name"),e.jsx(d,{name:"embedding_model_name",type:"string",children:e.jsx(n.p,{children:"Specified embedding model, corresponding to the model field(Optional)"})},"embedding_model_name"),e.jsx(d,{name:"score_threshold",type:"number",children:e.jsx(n.p,{children:"The similarity threshold for matching annotated replies. Only annotations with scores above this threshold will be recalled."})},"score_threshold")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.p,{children:"The provider and model name of the embedding model can be obtained through the following interface: v1/workspaces/current/models/model-types/text-embedding. For specific instructions, see: Maintain Knowledge Base via API. The Authorization used is the Dataset API Token."}),e.jsx(s,{title:"Request",tag:"POST",label:"/apps/annotation-reply/{action}",targetCode:`curl --location --request POST '${i.appDetail.api_base_url}/apps/annotation-reply/{action}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "score_threshold": 0.9,
  "embedding_provider_name": "zhipu",
  "embedding_model_name": "embedding_3"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "job_id": "b15c8f68-1cf4-4877-bf21-ed7cf2011802",
  "job_status": "waiting"
}
`})})}),e.jsx(n.p,{children:"This interface is executed asynchronously, so it will return a job_id. You can get the final execution result by querying the job status interface."})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotation-reply/{action}/status/{job_id}",method:"GET",title:"Query Initial Annotation Reply Settings Task Status",name:"#initial_annotation_reply_settings_task_status"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"action",type:"string",children:e.jsx(n.p,{children:"Action, can only be 'enable' or 'disable', must be the same as the action in the initial annotation reply settings interface"})},"action"),e.jsx(d,{name:"job_id",type:"string",children:e.jsx(n.p,{children:"Job ID, obtained from the initial annotation reply settings interface"})},"job_id")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/apps/annotations",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/apps/annotation-reply/{action}/status/{job_id}' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "job_id": "b15c8f68-1cf4-4877-bf21-ed7cf2011802",
  "job_status": "waiting",
  "error_msg": ""
}
`})})})]})]})]})}function Fn(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(ue,{...i})}):ue(i)}function pe(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"高度なチャットアプリ API"}),`
`,e.jsx(n.p,{children:"チャットアプリケーションはセッションの持続性をサポートしており、以前のチャット履歴を応答のコンテキストとして使用できます。これは、チャットボットやカスタマーサービス AI などに適用できます。"}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"ベース URL"}),e.jsx(s,{title:"コード",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"認証"}),e.jsxs(n.p,{children:["サービス API は ",e.jsx(n.code,{children:"API-Key"}),` 認証を使用します。
`,e.jsx("i",{children:e.jsx(n.strong,{children:"API キーはサーバー側に保存し、クライアント側で共有または保存しないことを強くお勧めします。API キーの漏洩は深刻な結果を招く可能性があります。"})})]}),e.jsxs(n.p,{children:["すべての API リクエストには、以下のように ",e.jsx(n.code,{children:"Authorization"}),"HTTP ヘッダーに API キーを含めてください："]}),e.jsx(s,{title:"コード",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages",method:"POST",title:"チャットメッセージを送信",name:"#Send-Chat-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"チャットアプリケーションにリクエストを送信します。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"query",type:"string",children:e.jsx(n.p,{children:"ユーザー入力/質問内容"})},"query"),e.jsx(d,{name:"inputs",type:"object",children:e.jsxs(n.p,{children:[`アプリによって定義されたさまざまな変数値の入力を許可します。
`,e.jsx(n.code,{children:"inputs"}),`パラメータには複数のキー/値ペアが含まれ、各キーは特定の変数に対応し、各値はその変数の特定の値です。
変数がファイルタイプの場合、以下の`,e.jsx(n.code,{children:"files"}),`で説明されているキーを持つオブジェクトを指定します。
デフォルト`,e.jsx(n.code,{children:"{}"})]})},"inputs"),e.jsxs(d,{name:"response_mode",type:"string",children:[e.jsx(n.p,{children:"応答の返却モードを指定します。サポートされているモード："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," ストリーミングモード（推奨）、SSE（",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"サーバー送信イベント"}),"）を通じてタイプライターのような出力を実装します。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` ブロッキングモード、実行完了後に結果を返します。（プロセスが長い場合、リクエストが中断される可能性があります）
Cloudflareの制限により、リクエストは100秒後に返答なしで中断されます。`]}),`
`]})]},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーの身元を定義するために使用され、統計のために使用されます。
アプリケーション内で開発者によって一意に定義されるべきです。サービス API は WebApp によって作成された会話を共有しません。`})},"user"),e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"会話ID、以前のチャット記録に基づいて会話を続けるには、以前のメッセージのconversation_idを渡す必要があります。"})},"conversation_id"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"ファイルリスト、モデルが Vision/Video 機能をサポートしている場合に限り、ファイルをテキスト理解および質問応答に組み合わせて入力するのに適しています。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) サポートされるタイプ：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," サポートされるタイプには以下が含まれます：'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," サポートされるタイプには以下が含まれます：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," サポートされるタイプには以下が含まれます：'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," サポートされるタイプには以下が含まれます：'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," サポートされるタイプには以下が含まれます：その他のファイルタイプ"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) 転送方法:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": ファイルのURL。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": ファイルをアップロード。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," ファイルのURL。（転送方法が ",e.jsx(n.code,{children:"remote_url"})," の場合のみ）。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," アップロードされたファイルID。（転送方法が ",e.jsx(n.code,{children:"local_file"})," の場合のみ）。"]}),`
`]})]},"files"),e.jsx(d,{name:"auto_generate_name",type:"bool",children:e.jsxs(n.p,{children:["タイトルを自動生成、デフォルトは",e.jsx(n.code,{children:"true"}),`。
`,e.jsx(n.code,{children:"false"}),"に設定すると、会話のリネームAPIを呼び出し、",e.jsx(n.code,{children:"auto_generate"}),"を",e.jsx(n.code,{children:"true"}),"に設定することで非同期タイトル生成を実現できます。"]})},"auto_generate_name"),e.jsx(d,{name:"workflow_id",type:"string",children:e.jsx(n.p,{children:"（オプション）ワークフローID、特定のバージョンを指定するために使用、提供されない場合はデフォルトの公開バージョンを使用。"})},"workflow_id"),e.jsxs(d,{name:"trace_id",type:"string",children:[e.jsxs(n.p,{children:["（オプション）トレースID。既存の業務システムのトレースコンポーネントと連携し、エンドツーエンドの分散トレーシングを実現するために使用します。指定がない場合、システムが自動的に trace_id を生成します。以下の3つの方法で渡すことができ、優先順位は次のとおりです：",e.jsx("br",{})]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["Header：HTTPヘッダー ",e.jsx("code",{children:"X-Trace-Id"})," で渡す（最優先）。",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["クエリパラメータ：URLクエリパラメータ ",e.jsx("code",{children:"trace_id"})," で渡す。",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["リクエストボディ：リクエストボディの ",e.jsx("code",{children:"trace_id"})," フィールドで渡す（本フィールド）。",e.jsx("br",{})]}),`
`]})]},"trace_id")]}),e.jsx(n.h3,{children:"応答"}),e.jsx(n.p,{children:`response_modeがブロッキングの場合、CompletionResponseオブジェクトを返します。
response_modeがストリーミングの場合、ChunkCompletionResponseストリームを返します。`}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["完全なアプリ結果を返します。",e.jsx(n.code,{children:"Content-Type"}),"は",e.jsx(n.code,{children:"application/json"}),"です。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) イベントタイプ、固定で ",e.jsx(n.code,{children:"message"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ユニークID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) アプリモード、",e.jsx(n.code,{children:"chat"}),"として固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 完全な応答内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) モデル使用情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用と帰属リスト"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) メッセージ作成タイムスタンプ、例：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["アプリによって出力されたストリームチャンクを返します。",e.jsx(n.code,{children:"Content-Type"}),"は",e.jsx(n.code,{children:"text/event-stream"}),`です。
各ストリーミングチャンクは`,e.jsx(n.code,{children:"data:"}),"で始まり、2つの改行文字",e.jsx(n.code,{children:"\\n\\n"}),"で区切られます。以下のように表示されます："]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsxs(n.p,{children:["ストリーミングチャンクの構造は",e.jsx(n.code,{children:"event"}),"に応じて異なります："]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLMがテキストチャンクイベントを返します。つまり、完全なテキストがチャンク形式で出力されます。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLMが返したテキストチャンク内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_file"})," メッセージファイルイベント、ツールによって新しいファイルが作成されました",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ファイル一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"}),' (string) ファイルタイプ、現在は"image"のみ許可']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) 所属、ここでは'assistant'のみ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) ファイルのリモートURL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"}),"  (string) 会話ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," メッセージ終了イベント、このイベントを受信するとストリーミングが終了したことを意味します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) モデル使用情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用と帰属リスト"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTSオーディオストリームイベント、つまり音声合成出力。内容はMp3形式のオーディオブロックで、base64文字列としてエンコードされています。再生時には、base64をデコードしてプレーヤーに入力するだけです。（このメッセージは自動再生が有効な場合にのみ利用可能）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のストップ応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 音声合成後のオーディオ、base64テキストコンテンツとしてエンコードされており、再生時にはbase64をデコードしてプレーヤーに入力するだけです"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTSオーディオストリーム終了イベント、このイベントを受信するとオーディオストリームが終了したことを示します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のストップ応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 終了イベントにはオーディオがないため、これは空の文字列です"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` メッセージ内容置換イベント。
出力内容のモデレーションが有効な場合、内容がフラグ付けされると、このイベントを通じてメッセージ内容がプリセットの返信に置き換えられます。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 置換内容（すべてのLLM返信テキストを直接置き換えます）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," ワークフローが実行を開始",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"workflow_started"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連ワークフローのID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," ノード実行が開始",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"node_started"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ノードのID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) ノードのタイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) ノードの名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 実行シーケンス番号、トレースノードシーケンスを表示するために使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) オプションのプレフィックスノードID、キャンバス表示実行パスに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ノードで使用されるすべての前のノード変数の内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始のタイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," ノード実行が終了、成功または失敗は同じイベント内で異なる状態で示されます",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"node_finished"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ノードのID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) ノードのタイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) ノードの名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 実行シーケンス番号、トレースノードシーケンスを表示するために使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) オプションのプレフィックスノードID、キャンバス表示実行パスに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ノードで使用されるすべての前のノード変数の内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) オプションのノードプロセスデータ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) オプションの出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行の状態、",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプションのエラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) オプションの使用される合計秒数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) オプションの使用されるトークン数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) オプションの合計コスト"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) オプション、例：",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始のタイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," ワークフロー実行が終了、成功または失敗は同じイベント内で異なる状態で示されます",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"workflow_finished"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行のID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連ワークフローのID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行の状態、",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) オプションの出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプションのエラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) オプションの使用される合計秒数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) オプションの使用されるトークン数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) デフォルト0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 終了時間"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
ストリーミングプロセス中に発生する例外はストリームイベントの形式で出力され、エラーイベントを受信するとストリームが終了します。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTPステータスコード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) エラーコード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) エラーメッセージ"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 接続を維持するために10秒ごとにpingイベントが発生します。"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404, 会話が存在しません"}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", 異常なパラメータ入力"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", アプリ構成が利用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", 利用可能なモデル資格情報構成がありません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", モデル呼び出しクォータが不足しています"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", 現在のモデルが利用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_not_found"}),", 指定されたワークフローバージョンが見つかりません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"draft_workflow_error"}),", ドラフトワークフローバージョンは使用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_id_format_error"}),", ワークフローID形式エラー、UUID形式が必要です"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"completion_request_error"}),", テキスト生成に失敗しました"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/chat-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "query": "What are the specs of the iPhone 13 Pro Max?",
  "response_mode": "streaming",
  "conversation_id": "",
  "user": "abc-123",
  "files": [
      {
          "type": "image",
          "transfer_method": "remote_url",
          "url": "https://cloud.gofy.ai/logo/logo-site.png"
      }
  ]
}'`}),e.jsx(n.h3,{children:"ブロッキングモード"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "event": "message",
    "task_id": "c3800678-a077-43df-a102-53f23ed20b88",
    "id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "message_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
    "mode": "chat",
    "answer": "iPhone 13 Pro Maxの仕様は次のとおりです:...",
    "metadata": {
        "usage": {
            "prompt_tokens": 1033,
            "prompt_unit_price": "0.001",
            "prompt_price_unit": "0.001",
            "prompt_price": "0.0010330",
            "completion_tokens": 128,
            "completion_unit_price": "0.002",
            "completion_price_unit": "0.001",
            "completion_price": "0.0002560",
            "total_tokens": 1161,
            "total_price": "0.0012890",
            "currency": "USD",
            "latency": 0.7682376249867957
        },
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ]
    },
    "created_at": 1705407629
}
`})})}),e.jsx(n.h3,{children:"ストリーミングモード"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " I", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": "'m", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " glad", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " to", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " meet", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " you", "created_at": 1679586595}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"ファイルアップロード",name:"#file-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:`メッセージ送信時に使用するファイルをアップロードし、画像とテキストのマルチモーダル理解を可能にします。
アプリケーションでサポートされている形式をサポートします。
アップロードされたファイルは現在のエンドユーザーのみが使用できます。`}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.p,{children:["このインターフェースは",e.jsx(n.code,{children:"multipart/form-data"}),"リクエストを必要とします。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file"}),` (File) 必須
アップロードするファイル。`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子、開発者のルールによって定義され、アプリケーション内で一意でなければなりません。サービス API は WebApp によって作成された会話を共有しません。`]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsx(n.p,{children:"アップロードが成功すると、サーバーはファイルの ID と関連情報を返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) ファイル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) ファイルサイズ（バイト）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) ファイル拡張子"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) ファイルの MIME タイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) エンドユーザーID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"no_file_uploaded"}),", ファイルが提供されなければなりません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"too_many_files"}),", 現在は 1 つのファイルのみ受け付けます"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_preview"}),", ファイルはプレビューをサポートしていません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_estimate"}),", ファイルは推定をサポートしていません"]}),`
`,e.jsxs(n.li,{children:["413, ",e.jsx(n.code,{children:"file_too_large"}),", ファイルが大きすぎます"]}),`
`,e.jsxs(n.li,{children:["415, ",e.jsx(n.code,{children:"unsupported_file_type"}),", サポートされていない拡張子、現在はドキュメントファイルのみ受け付けます"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_connection_failed"}),", S3 サービスに接続できません"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_permission_denied"}),", S3 にファイルをアップロードする権限がありません"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_file_too_large"}),", ファイルが S3 のサイズ制限を超えています"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"エンドユーザーを取得",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"エンドユーザー ID からエンドユーザー情報を取得します。"}),e.jsxs(n.p,{children:["他の API がエンドユーザー ID（例：ファイルアップロードの ",e.jsx(n.code,{children:"created_by"}),"）を返す場合に利用できます。"]}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) 必須
エンドユーザー ID。`]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"EndUser オブジェクトを返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) テナント ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) アプリ ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) エンドユーザー種別"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) 外部ユーザー ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) 匿名ユーザーかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) セッション ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 日時"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 日時"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"end_user_not_found"}),", エンドユーザーが見つかりません"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"レスポンス例"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"ファイルプレビュー",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アップロードされたファイルをプレビューまたはダウンロードします。このエンドポイントを使用すると、以前にファイルアップロード API でアップロードされたファイルにアクセスできます。"}),e.jsx("i",{children:"ファイルは、リクエストしているアプリケーションのメッセージ範囲内にある場合のみアクセス可能です。"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) 必須
プレビューするファイルの一意識別子。ファイルアップロード API レスポンスから取得します。`]}),`
`]}),e.jsx(n.h3,{children:"クエリパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) オプション
ファイルを添付ファイルとして強制ダウンロードするかどうか。デフォルトは `,e.jsx(n.code,{children:"false"}),"（ブラウザでプレビュー）。"]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"ブラウザ表示またはダウンロード用の適切なヘッダー付きでファイル内容を返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," ファイル MIME タイプに基づいて設定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," ファイルサイズ（バイト、利用可能な場合）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"})," ",e.jsx(n.code,{children:"as_attachment=true"}),' の場合は "attachment" に設定']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," パフォーマンス向上のためのキャッシュヘッダー"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' 音声/動画ファイルの場合は "bytes" に設定']}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", パラメータ入力異常"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", ファイルアクセス拒否またはファイルが現在のアプリケーションに属していません"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", ファイルが見つからないか削除されています"]}),`
`,e.jsx(n.li,{children:"500, サーバー内部エラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"添付ファイルとしてダウンロード"}),e.jsx(s,{title:"ダウンロードリクエスト",tag:"GET",label:"/files/:file_id/preview?as_attachment=true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"レスポンスヘッダー例"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"ダウンロードレスポンスヘッダー"}),e.jsx(s,{title:"Download Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages/:task_id/stop",method:"POST",title:"生成を停止",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ストリーミングモードでのみサポートされています。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、ストリーミングチャンクの返り値から取得できます"]}),`
`]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子、エンドユーザーの身元を定義するために使用され、送信メッセージインターフェースで渡されたユーザーと一致している必要があります。サービス API は WebApp によって作成された会話を共有しません。`]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"POST",label:"/chat-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{
  "user": "abc-123"
}'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"メッセージフィードバック",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"エンドユーザーはフィードバックメッセージを提供でき、アプリケーション開発者が期待される出力を最適化するのを支援します。"}),e.jsx(n.h3,{children:"パス"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"メッセージID"})},"message_id")}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsxs(n.p,{children:["アップボートは",e.jsx(n.code,{children:"like"}),"、ダウンボートは",e.jsx(n.code,{children:"dislike"}),"、アップボートの取り消しは",e.jsx(n.code,{children:"null"})]})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者のルールによって定義され、アプリケーション内で一意でなければなりません。"})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"メッセージのフィードバックです。"})},"content")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"アプリのメッセージの「いいね」とフィードバックを取得",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アプリのエンドユーザーからのフィードバックや「いいね」を取得します。"}),e.jsx(n.h3,{children:"クエリ"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（任意）ページ番号。デフォルト値：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（任意）1ページあたりの件数。デフォルト値：20"})},"limit")}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (リスト) このアプリの「いいね」とフィードバックの一覧を返します。"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`    {
    "data": [
        {
            "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
            "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
            "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
            "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
            "rating": "like",
            "content": "message feedback information-3",
            "from_source": "user",
            "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
            "from_account_id": null,
            "created_at": "2025-04-24T09:24:38",
            "updated_at": "2025-04-24T09:24:38"
        }
    ]
    }
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/{message_id}/suggested",method:"GET",title:"次の推奨質問",name:"#suggested"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"現在のメッセージに対する次の質問の提案を取得します"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"メッセージID"})},"message_id")}),e.jsx(n.h3,{children:"クエリ"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーの身元を定義するために使用され、統計のために使用されます。
アプリケーション内で開発者によって一意に定義されるべきです。`})},"user")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/messages/{message_id}/suggested",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/messages/{message_id}/suggested?user=abc-123' \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success",
  "data": [
        "a",
        "b",
        "c"
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages",method:"GET",title:"会話履歴メッセージを取得",name:"#messages"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:["スクロールロード形式で履歴チャット記録を返し、最初のページは最新の",e.jsx(n.code,{children:"{limit}"}),"メッセージを返します。つまり、逆順です。"]}),e.jsx(n.h3,{children:"クエリ"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"会話ID"})},"conversation_id"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーの身元を定義するために使用され、統計のために使用されます。
アプリケーション内で開発者によって一意に定義されるべきです。`})},"user"),e.jsx(d,{name:"first_id",type:"string",children:e.jsx(n.p,{children:"現在のページの最初のチャット記録のID、デフォルトはnullです。"})},"first_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"1回のリクエストで返すチャット履歴メッセージの数、デフォルトは20です。"})},"limit")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) メッセージリスト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) メッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ユーザー入力パラメータ。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"})," (string) ユーザー入力/質問内容。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[object]) メッセージファイル",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) ファイルタイプ、画像の場合はimage"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) ファイルプレビューURL、ファイルアクセスにはファイルプレビューAPI（",e.jsx(n.code,{children:"/files/{file_id}/preview"}),"）を使用してください"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) 所属、userまたはassistant"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 応答メッセージ内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"feedback"})," (object) フィードバック情報",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"rating"})," (string) アップボートは",e.jsx(n.code,{children:"like"})," / ダウンボートは",e.jsx(n.code,{children:"dislike"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用と帰属リスト"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) 次のページがあるかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返された項目数、入力がシステム制限を超える場合、システム制限数を返します"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/messages",targetCode:`curl -X GET '${i.appDetail.api_base_url}/messages?user=abc-123&conversation_id={conversation_id}'
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
        "id": "a076a87f-31e5-48dc-b452-0061adbbc922",
        "conversation_id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
        "inputs": {
            "name": "gofy"
        },
        "query": "iphone 13 pro",
        "answer": "iPhone 13 Proは2021年9月24日に発売され、6.1インチのディスプレイと1170 x 2532の解像度を備えています。Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)プロセッサ、6 GBのRAMを搭載し、128 GB、256 GB、512 GB、1 TBのストレージオプションを提供します。カメラは12 MP、バッテリー容量は3095 mAhで、iOS 15を搭載しています。",
        "message_files": [],
        "feedback": null,
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ],
        "created_at": 1705569239,
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations",method:"GET",title:"会話を取得",name:"#conversations"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"現在のユーザーの会話リストを取得し、デフォルトで最新の 20 件を返します。"}),e.jsx(n.h3,{children:"クエリ"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーの身元を定義するために使用され、統計のために使用されます。
アプリケーション内で開発者によって一意に定義されるべきです。`})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"(Optional)現在のページの最後の記録のID、デフォルトはnullです。"})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"(Optional)1回のリクエストで返す記録の数、デフォルトは最新の20件です。最大100、最小1。"})},"limit"),e.jsxs(d,{name:"sort_by",type:"string",children:[e.jsx(n.p,{children:"(Optional)ソートフィールド、デフォルト：-updated_at（更新時間で降順にソート）"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"利用可能な値：created_at, -created_at, updated_at, -updated_at"}),`
`,e.jsx(n.li,{children:'フィールドの前の記号は順序または逆順を表し、"-"は逆順を表します。'}),`
`]})]},"sort_by")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 会話のリスト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 会話名、デフォルトではLLMによって生成されます。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ユーザー入力パラメータ。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) 紹介"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) 更新タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返されたエントリ数、入力がシステム制限を超える場合、システム制限数が返されます"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/conversations",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations?user=abc-123&last_id=&limit=20' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
      "id": "10799fb8-64f7-4296-bbf7-b42bfbe0ae54",
      "name": "新しいチャット",
      "inputs": {
          "book": "book",
          "myName": "Lucy"
      },
      "status": "normal",
      "created_at": 1679667915,
      "updated_at": 1679667915
    },
    {
      "id": "hSIhXBhNe8X1d8Et"
      // ...
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id",method:"DELETE",title:"会話を削除",name:"#delete"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"会話を削除します。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者によって定義され、アプリケーション内で一意であることを保証しなければなりません。"})},"user")}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"DELETE",label:"/conversations/:conversation_id",targetCode:`curl -X DELETE '${i.appDetail.api_base_url}/conversations/{conversation_id}' \\
--header 'Content-Type: application/json' \\
--header 'Accept: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data '{
  "user": "abc-123"
}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/name",method:"POST",title:"会話の名前を変更",name:"#rename"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"リクエストボディ"}),e.jsx(n.p,{children:"セッションの名前を変更します。セッション名は、複数のセッションをサポートするクライアントでの表示に使用されます。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`]}),e.jsxs(t,{children:[e.jsx(d,{name:"name",type:"string",children:e.jsxs(n.p,{children:["(Optional)会話の名前。",e.jsx(n.code,{children:"auto_generate"}),"が",e.jsx(n.code,{children:"true"}),"に設定されている場合、このパラメータは省略できます。"]})},"name"),e.jsx(d,{name:"auto_generate",type:"bool",children:e.jsxs(n.p,{children:["(Optional)タイトルを自動生成、デフォルトは",e.jsx(n.code,{children:"false"})]})},"auto_generate"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者によって定義され、アプリケーション内で一意であることを保証しなければなりません。"})},"user")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 会話名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ユーザー入力パラメータ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 会話状態"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) 紹介"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) 更新タイムスタンプ、例：1705395332"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/conversations/:conversation_id/name",targetCode:`curl -X POST '${i.appDetail.api_base_url}/conversations/{conversation_id}/name' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "name": "",
  "auto_generate": true,
  "user": "abc-123"
}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
    "name": "チャット vs AI",
    "inputs": {},
    "status": "normal",
    "introduction": "",
    "created_at": 1705569238,
    "updated_at": 1705569238
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables",method:"GET",title:"会話変数の取得",name:"#conversation-variables"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"特定の会話から変数を取得します。このエンドポイントは、会話中に取得された構造化データを抽出するのに役立ちます。"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsx(t,{children:e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"変数を取得する会話のID。"})},"conversation_id")}),e.jsx(n.h3,{children:"クエリパラメータ"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子。開発者によって定義されたルールに従い、アプリケーション内で一意である必要があります。"})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"(Optional)現在のページの最後の記録のID、デフォルトはnullです。"})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"(Optional)1回のリクエストで返す記録の数、デフォルトは最新の20件です。最大100、最小1。"})},"limit")]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) ページごとのアイテム数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) さらにアイテムがあるかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 変数のリスト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 変数 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 変数名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) 変数タイプ（文字列、数値、真偽値など）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (string) 変数値"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 変数の説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) 最終更新タイムスタンプ"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", 会話が見つかりません"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/conversations/:conversation_id/variables",debug:"true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"変数名フィルター付きリクエスト",language:"bash",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123&variable_name=customer_name' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 100,
  "has_more": false,
  "data": [
    {
      "id": "variable-uuid-1",
      "name": "customer_name",
      "value_type": "string",
      "value": "John Doe",
      "description": "会話から抽出された顧客名",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    },
    {
      "id": "variable-uuid-2",
      "name": "order_details",
      "value_type": "json",
      "value": "{\\"product\\":\\"Widget\\",\\"quantity\\":5,\\"price\\":19.99}",
      "description": "顧客の注文詳細",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables/:variable_id",method:"PUT",title:"会話変数の更新",name:"#update-conversation-variable"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"特定の会話変数の値を更新します。このエンドポイントは、名前、型、説明を保持しながら、会話中にキャプチャされた変数の値を変更することを可能にします。"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"更新する変数を含む会話のID。"})},"conversation_id"),e.jsx(d,{name:"variable_id",type:"string",children:e.jsx(n.p,{children:"更新する変数のID。"})},"variable_id")]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"value",type:"any",children:e.jsx(n.p,{children:"変数の新しい値。変数の期待される型（文字列、数値、オブジェクトなど）と一致する必要があります。"})},"value"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子。開発者によって定義されたルールに従い、アプリケーション内で一意である必要があります。"})},"user")]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"以下を含む更新された変数オブジェクトを返します："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 変数名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) 変数型（文字列、数値、オブジェクトなど）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (any) 更新された変数値"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 変数の説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) 最終更新タイムスタンプ"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"Type mismatch: variable expects {expected_type}, but got {actual_type} type"}),", 値の型が変数の期待される型と一致しません"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", 会話が見つかりません"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_variable_not_exists"}),", 変数が見つかりません"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"PUT",label:"/conversations/:conversation_id/variables/:variable_id",targetCode:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": "Updated Value",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"異なる値型での更新",targetCode:[{title:"文字列値",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": "新しい文字列値",
  "user": "abc-123"
}'`},{title:"数値",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": 42,
  "user": "abc-123"
}'`},{title:"オブジェクト値",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": {"product": "Widget", "quantity": 10, "price": 29.99},
  "user": "abc-123"
}'`}]}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "variable-uuid-1",
  "name": "customer_name",
  "value_type": "string",
  "value": "Updated Value",
  "description": "会話から抽出された顧客名",
  "created_at": 1650000000000,
  "updated_at": 1650000001000
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/audio-to-text",method:"POST",title:"音声からテキストへ",name:"#audio-to-text"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"このエンドポイントは multipart/form-data リクエストを必要とします。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsxs(n.p,{children:[`オーディオファイル。
サポートされている形式：`,e.jsx(n.code,{children:"['mp3', 'mp4', 'mpeg', 'mpga', 'm4a', 'wav', 'webm']"}),`
ファイルサイズ制限：15MB`]})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者のルールによって定義され、アプリケーション内で一意でなければなりません。"})},"user")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) 出力テキスト"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/audio-to-text",targetCode:`curl -X POST '${i.appDetail.api_base_url}/audio-to-text' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=audio/[mp3|mp4|mpeg|mpga|m4a|wav|webm]'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "text": ""
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"テキストから音声へ",name:"#text-to-audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"テキストを音声に変換します。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"Gofyによって生成されたテキストメッセージの場合、生成されたメッセージIDを直接渡します。バックエンドはメッセージIDを使用して対応する内容を検索し、音声情報を直接合成します。message_idとtextが同時に提供される場合、message_idが優先されます。"})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"音声生成コンテンツ。"})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者によって定義され、アプリ内で一意であることを保証しなければなりません。"})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/text-to-audio",targetCode:`curl -o text-to-audio.mp3 -X POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290",
  "text": "Hello Gofy",
  "user": "abc-123",
}'`}),e.jsx(s,{title:"ヘッダー",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"アプリケーションの基本情報を取得",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"このアプリケーションの基本情報を取得するために使用されます"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) アプリケーションの名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) アプリケーションの説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) アプリケーションのタグ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) アプリケーションのモード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"author_name"})," (string) 作者の名前"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "advanced-chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"アプリケーションのパラメータ情報を取得",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ページに入る際に、機能、入力パラメータ名、タイプ、デフォルト値などの情報を取得するために使用されます。"}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) 開始の挨拶"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) 開始時の推奨質問のリスト"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) 答えを有効にした後の質問を提案します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) 音声からテキストへ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text_to_speech"})," (object) テキストから音声へ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"voice"})," (string) 音声タイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"language"})," (string) 言語"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"autoPlay"})," (string) 自動再生",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"}),"  有効"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"disabled"})," 無効"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) 引用と帰属",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) 注釈返信",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) ユーザー入力フォームの設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) テキスト入力コントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) 段落テキスト入力コントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) ドロップダウンコントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) オプション値"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) ファイルアップロード設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) ドキュメント設定
現在サポートされているドキュメントタイプ：`,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) ドキュメント数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) 画像設定
現在サポートされている画像タイプ：`,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 画像数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) オーディオ設定
現在サポートされているオーディオタイプ：`,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) オーディオ数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) ビデオ設定
現在サポートされているビデオタイプ：`,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) ビデオ数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) カスタム設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) カスタム数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) システムパラメータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) ドキュメントアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) 画像ファイルアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) オーディオファイルアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) ビデオファイルアップロードサイズ制限（MB）"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/parameters",targetCode:`curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "opening_statement": "こんにちは！",
  "suggested_questions_after_answer": {
      "enabled": true
  },
  "speech_to_text": {
      "enabled": true
  },
  "text_to_speech": {
      "enabled": true,
      "voice": "sambert-zhinan-v1",
      "language": "zh-Hans",
      "autoPlay": "disabled"
  },
  "retriever_resource": {
      "enabled": true
  },
  "annotation_reply": {
      "enabled": true
  },
  "user_input_form": [
      {
          "paragraph": {
              "label": "クエリ",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/meta",method:"GET",title:"アプリケーションのメタ情報を取得",name:"#meta"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"このアプリケーションのツールのアイコンを取得するために使用されます"}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_icons"}),"(object[string]) ツールアイコン",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_name"})," (string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (object|string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["(object) アイコンオブジェクト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"background"})," (string) 背景色（16 進数形式）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"content"}),"(string) 絵文字"]}),`
`]}),`
`]}),`
`,e.jsx(n.li,{children:"(string) アイコンの URL"}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/meta",targetCode:`curl -X GET '${i.appDetail.api_base_url}/meta' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "tool_icons": {
    "dalle2": "https://cloud.gofy.ai/console/api/workspaces/current/tool-provider/builtin/dalle/icon",
    "api_tool": {
      "background": "#252525",
      "content": "\\ud83d\\ude01"
    }
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"アプリのWebApp設定を取得",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アプリの WebApp 設定を取得するために使用します。"}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp 名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) チャットの色テーマ、16 進数形式"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) チャットの色テーマを反転するかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) アイコンタイプ、",e.jsx(n.code,{children:"emoji"}),"-絵文字、",e.jsx(n.code,{children:"image"}),"-画像"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) アイコン。",e.jsx(n.code,{children:"emoji"}),"タイプの場合は絵文字、",e.jsx(n.code,{children:"image"}),"タイプの場合は画像 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) 16 進数形式の背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) アイコンの URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) 著作権情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) プライバシーポリシーのリンク"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) カスタム免責事項"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) デフォルト言語"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) ワークフローの詳細を表示するかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) WebApp のアイコンをチャット内の🤖に置き換えるかどうか"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{})]})}function Vn(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(pe,{...i})}):pe(i)}function _e(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"工作流编排对话型应用 API"}),`
`,e.jsx(n.p,{children:"对话应用支持会话持久化，可将之前的聊天记录作为上下文进行回答，可适用于聊天/客服 AI 等。"}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"基础 URL"}),e.jsx(s,{title:"Code",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"鉴权"}),e.jsxs(n.p,{children:["Service API 使用 ",e.jsx(n.code,{children:"API-Key"}),` 进行鉴权。
`,e.jsx("i",{children:e.jsxs(n.strong,{children:["强烈建议开发者把 ",e.jsx(n.code,{children:"API-Key"})," 放在后端存储，而非分享或者放在客户端存储，以免 ",e.jsx(n.code,{children:"API-Key"})," 泄露，导致财产损失。"]})}),`
所有 API 请求都应在 `,e.jsx(n.strong,{children:e.jsx(n.code,{children:"Authorization"})})," HTTP Header 中包含您的 ",e.jsx(n.code,{children:"API-Key"}),"，如下所示："]}),e.jsx(s,{title:"Code",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages",method:"POST",title:"发送对话消息",name:"#Create-Chat-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"创建会话消息。"}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"query",type:"string",children:e.jsx(n.p,{children:"用户输入/提问内容。"})},"query"),e.jsx(d,{name:"inputs",type:"object",children:e.jsxs(n.p,{children:[`允许传入 App 定义的各变量值。
inputs 参数包含了多组键值对（Key/Value pairs），每组的键对应一个特定变量，每组的值则是该变量的具体值。
如果变量是文件类型，请指定一个包含以下 `,e.jsx(n.code,{children:"files"}),` 中所述键的对象。
默认 `,e.jsx(n.code,{children:"{}"})]})},"inputs"),e.jsx(d,{name:"response_mode",type:"string",children:e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," 流式模式（推荐）。基于 SSE（",e.jsx(n.strong,{children:e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"})}),"）实现类似打字机输出方式的流式返回。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` 阻塞模式，等待执行完毕后返回结果。（请求若流程较长可能会被中断）。
`,e.jsx("i",{children:"由于 Cloudflare 限制，请求会在 100 秒超时无返回后中断。"})]}),`
`]})},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`用户标识，用于定义终端用户的身份，方便检索、统计。
由开发者定义规则，需保证用户标识在应用内唯一。服务 API 不会共享 WebApp 创建的对话。`})},"user"),e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"（选填）会话 ID，需要基于之前的聊天记录继续对话，必须传之前消息的 conversation_id。"})},"conversation_id"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"文件列表，适用于传入文件结合文本理解并回答问题，仅当模型支持 Vision/Video 能力时可用。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 支持类型：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," 具体类型包含：'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," 具体类型包含：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," 具体类型包含：'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," 具体类型包含：'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," 具体类型包含：其他文件类型"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string)  传递方式:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": 文件地址。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": 上传文件。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," 文件地址。（仅当传递方式为 ",e.jsx(n.code,{children:"remote_url"})," 时）。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," 上传文件 ID。（仅当传递方式为 ",e.jsx(n.code,{children:"local_file "}),"时）。"]}),`
`]})]},"files"),e.jsx(d,{name:"auto_generate_name",type:"bool",children:e.jsxs(n.p,{children:["（选填）自动生成标题，默认 ",e.jsx(n.code,{children:"true"}),"。 若设置为 ",e.jsx(n.code,{children:"false"}),"，则可通过调用会话重命名接口并设置 ",e.jsx(n.code,{children:"auto_generate"})," 为 ",e.jsx(n.code,{children:"true"})," 实现异步生成标题。"]})},"auto_generate_name"),e.jsx(d,{name:"workflow_id",type:"string",children:e.jsx(n.p,{children:"（选填）工作流ID，用于指定特定版本，如果不提供则使用默认的已发布版本。"})},"workflow_id"),e.jsxs(d,{name:"trace_id",type:"string",children:[e.jsxs(n.p,{children:["（选填）链路追踪ID。适用于与业务系统已有的trace组件打通，实现端到端分布式追踪等场景。如果未指定，系统会自动生成",e.jsx("code",{children:"trace_id"}),"。支持以下三种方式传递，具体优先级依次为：",e.jsx("br",{})]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["Header：通过 HTTP Header ",e.jsx("code",{children:"X-Trace-Id"})," 传递，优先级最高。",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["Query 参数：通过 URL 查询参数 ",e.jsx("code",{children:"trace_id"})," 传递。",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["Request Body：通过请求体字段 ",e.jsx("code",{children:"trace_id"})," 传递（即本字段）。",e.jsx("br",{})]}),`
`]})]},"trace_id")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(t,{children:[e.jsxs(n.p,{children:["当 ",e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"blocking"}),` 时，返回 ChatCompletionResponse object。
当 `,e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"streaming"}),"时，返回 ChunkChatCompletionResponse object 流式序列。"]}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["返回完整的 App 结果，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"application/json"}),"。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 事件类型，固定为 ",e.jsx(n.code,{children:"message"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 唯一ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) App 模式，固定为 chat"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 完整回复内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) 模型用量信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用和归属分段列表"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 消息创建时间戳，如：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["返回 App 输出的流式块，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"text/event-stream"}),`。
每个流式块均为 data: 开头，块之间以 \\n\\n 即两个换行符分隔，如下所示：`]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsx(n.p,{children:"流式块中根据 event 不同，结构也不同："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLM 返回文本块事件，即：完整的文本以分块的方式输出。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLM 返回文本块内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_file"})," 文件事件，表示有新文件需要展示",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 文件唯一ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 文件类型，目前仅为image"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) 文件归属，user或assistant，该接口返回仅为 ",e.jsx(n.code,{children:"assistant"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) 文件访问地址"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"}),"  (string) 会话ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," 消息结束事件，收到此事件则代表流式返回结束。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) 模型用量信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用和归属分段列表"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS 音频流事件，即：语音合成输出。内容是Mp3格式的音频块，使用 base64 编码后的字符串，播放的时候直接解码即可。(开启自动播放才有此消息)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 语音合成之后的音频块使用 Base64 编码之后的文本内容，播放的时候直接 base64 解码送入播放器即可"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS 音频流结束事件，收到这个事件表示音频流返回结束。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 结束事件是没有音频的，所以这里是空字符串"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` 消息内容替换事件。
开启内容审查和审查输出内容时，若命中了审查条件，则会通过此事件替换消息内容为预设回复。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 替换内容（直接替换 LLM 所有回复文本）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," workflow 开始执行",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"workflow_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," node 开始执行",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"node_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) 节点 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) 节点类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) 节点名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 执行序号，用于展示 Tracing Node 顺序"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) 前置节点 ID，用于画布展示执行路径"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 节点中所有使用到的前置节点变量内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," node 执行结束，成功失败同一事件中不同状态",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"node_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) node 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) 节点 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 执行序号，用于展示 Tracing Node 顺序"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional 前置节点 ID，用于画布展示执行路径"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 节点中所有使用到的前置节点变量内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) Optional 节点过程数据"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional 输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 执行状态 ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional 耗时(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) optional 总使用 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) optional 总费用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) optional 货币，如 ",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," workflow 执行结束，成功失败同一事件中不同状态",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"workflow_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string)  执行状态 ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional 输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional 耗时(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional 总使用 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 总步数（冗余），默认 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 结束时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
流式输出过程中出现的异常会以 stream event 形式输出，收到异常事件后即结束。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTP 状态码"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) 错误码"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) 错误消息"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 每 10s 一次的 ping 事件，保持连接存活。"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404，对话不存在"}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"invalid_param"}),"，传入参数异常"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"app_unavailable"}),"，App 配置不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_not_initialize"}),"，无可用模型凭据配置"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_quota_exceeded"}),"，模型调用额度不足"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"model_currently_not_support"}),"，当前模型不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_not_found"}),"，指定的工作流版本未找到"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"draft_workflow_error"}),"，无法使用草稿工作流版本"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_id_format_error"}),"，工作流ID格式错误，需要UUID格式"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"completion_request_error"}),"，文本生成失败"]}),`
`,e.jsx(n.li,{children:"500，服务内部异常"}),`
`]})]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/chat-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "query": "What are the specs of the iPhone 13 Pro Max?",
  "response_mode": "streaming",
  "conversation_id": "",
  "user": "abc-123",
  "files": [
      {
          "type": "image",
          "transfer_method": "remote_url",
          "url": "https://cloud.gofy.ai/logo/logo-site.png"
      }
  ]
}'`}),e.jsx(n.h3,{children:"阻塞模式"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "event": "message",
    "task_id": "c3800678-a077-43df-a102-53f23ed20b88",
    "id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "message_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
    "mode": "chat",
    "answer": "iPhone 13 Pro Max specs are listed here:...",
    "metadata": {
        "usage": {
            "prompt_tokens": 1033,
            "prompt_unit_price": "0.001",
            "prompt_price_unit": "0.001",
            "prompt_price": "0.0010330",
            "completion_tokens": 128,
            "completion_unit_price": "0.002",
            "completion_price_unit": "0.001",
            "completion_price": "0.0002560",
            "total_tokens": 1161,
            "total_price": "0.0012890",
            "currency": "USD",
            "latency": 0.7682376249867957
        },
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ]
    },
    "created_at": 1705407629
}
`})})}),e.jsx(n.h3,{children:"流式模式"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " I", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": "'m", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " glad", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " to", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " meet", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " you", "created_at": 1679586595}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"上传文件",name:"#files-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:[`上传文件并在发送消息时使用，可实现图文多模态理解。
支持您的应用程序所支持的所有格式。
`,e.jsx("i",{children:"上传的文件仅供当前终端用户使用。"})]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["该接口需使用  ",e.jsx(n.code,{children:"multipart/form-data"})," 进行请求。"]}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsx(n.p,{children:"要上传的文件。"})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，用于定义终端用户的身份，必须和发送消息接口传入 user 保持一致。"})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"成功上传后，服务器会返回文件的 ID 和相关信息。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 文件名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) 文件大小（byte）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) 文件后缀"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) 文件 mime-type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) 上传人 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 上传时间"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"no_file_uploaded"}),"，必须提供文件"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"too_many_files"}),"，目前只接受一个文件"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"unsupported_preview"}),"，该文件不支持预览"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"unsupported_estimate"}),"，该文件不支持估算"]}),`
`,e.jsxs(n.li,{children:["413，",e.jsx(n.code,{children:"file_too_large"}),"，文件太大"]}),`
`,e.jsxs(n.li,{children:["415，",e.jsx(n.code,{children:"unsupported_file_type"}),"，不支持的扩展名，当前只接受文档类文件"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_connection_failed"}),"，无法连接到 S3 服务"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_permission_denied"}),"，无权限上传文件到 S3"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_file_too_large"}),"，文件超出 S3 大小限制"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": 123,
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"获取终端用户",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"通过终端用户 ID 获取终端用户信息。"}),e.jsxs(n.p,{children:["当其他 API 返回终端用户 ID（例如：上传文件接口返回的 ",e.jsx(n.code,{children:"created_by"}),"）时，可使用该接口查询对应的终端用户信息。"]}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) 必需
终端用户 ID。`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"返回 EndUser 对象。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) 工作空间（Tenant）ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) 应用 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 终端用户类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) 外部用户 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) 是否匿名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 时间"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404，",e.jsx(n.code,{children:"end_user_not_found"}),"，终端用户不存在"]}),`
`,e.jsx(n.li,{children:"500，内部服务器错误"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"文件预览",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"预览或下载已上传的文件。此端点允许您访问先前通过文件上传 API 上传的文件。"}),e.jsx("i",{children:"文件只能在属于请求应用程序的消息范围内访问。"}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) 必需
要预览的文件的唯一标识符，从文件上传 API 响应中获得。`]}),`
`]}),e.jsx(n.h3,{children:"查询参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) 可选
是否强制将文件作为附件下载。默认为 `,e.jsx(n.code,{children:"false"}),"（在浏览器中预览）。"]}),`
`]}),e.jsx(n.h3,{children:"响应"}),e.jsx(n.p,{children:"返回带有适当浏览器显示或下载标头的文件内容。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," 根据文件 MIME 类型设置"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," 文件大小（以字节为单位，如果可用）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"})," 如果 ",e.jsx(n.code,{children:"as_attachment=true"}),' 则设置为 "attachment"']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," 用于性能的缓存标头"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' 对于音频/视频文件设置为 "bytes"']}),`
`]}),e.jsx(n.h3,{children:"错误"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", 参数输入异常"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", 文件访问被拒绝或文件不属于当前应用程序"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", 文件未找到或已被删除"]}),`
`,e.jsx(n.li,{children:"500, 服务内部错误"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"请求示例"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"作为附件下载"}),e.jsx(s,{title:"Download Request",tag:"GET",label:"/files/:file_id/preview?as_attachment=true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"响应标头示例"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"文件下载响应标头"}),e.jsx(s,{title:"Download Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages/:task_id/stop",method:"POST",title:"停止响应",name:"#Stop"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"仅支持流式模式。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，可在流式返回 Chunk 中获取"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
用户标识，用于定义终端用户的身份，必须和发送消息接口传入 user 保持一致。API 无法访问 WebApp 创建的会话。`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"})," (string) 固定返回 success"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/chat-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"消息反馈（点赞）",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"消息终端用户反馈、点赞，方便应用开发者优化输出预期。"}),e.jsx(n.h3,{children:"Path Params"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"消息 ID"})},"message_id")}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsx(n.p,{children:"点赞 like, 点踩 dislike,  撤销点赞 null"})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"消息反馈的具体信息。"})},"content")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"})," (string) 固定返回 success"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"获取APP的消息点赞和反馈",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"获取应用的终端用户反馈、点赞。"}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（选填）分页，默认值：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（选填）每页数量，默认值：20"})},"limit")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (List) 返回该APP的点赞、反馈列表。"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`    {
    "data": [
        {
            "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
            "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
            "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
            "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
            "rating": "like",
            "content": "message feedback information-3",
            "from_source": "user",
            "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
            "from_account_id": null,
            "created_at": "2025-04-24T09:24:38",
            "updated_at": "2025-04-24T09:24:38"
        }
    ]
    }

`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/{message_id}/suggested",method:"GET",title:"获取下一轮建议问题列表",name:"#suggested"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"获取下一轮建议问题列表。"}),e.jsx(n.h3,{children:"Path Params"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"Message ID"})},"message_id")}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/messages/{message_id}/suggested",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/messages/{message_id}/suggested?user=abc-123' \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success",
  "data": [
        "a",
        "b",
        "c"
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages",method:"GET",title:"获取会话历史消息",name:"#messages"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:["滚动加载形式返回历史聊天记录，第一页返回最新  ",e.jsx(n.code,{children:"limit"})," 条，即：倒序返回。"]}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"会话 ID"})},"conversation_id"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user"),e.jsx(d,{name:"first_id",type:"string",children:e.jsx(n.p,{children:"当前页第一条聊天记录的 ID，默认 null"})},"first_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"一次请求返回多少条聊天记录，默认 20 条。"})},"limit")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object])  消息列表",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"}),"  (string) 消息 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string)  会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 用户输入参数。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"}),"  (string) 用户输入 / 提问内容。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[object]) 消息文件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 文件类型，image 图片"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) 文件预览地址，使用文件预览 API (",e.jsx(n.code,{children:"/files/{file_id}/preview"}),") 访问文件"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) 文件归属方，user 或 assistant"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string)  回答消息内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"}),"  (timestamp) 创建时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"feedback"})," (object) 反馈信息",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"rating"})," (string) 点赞 like / 点踩 dislike"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用和归属分段列表"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) 是否存在下一页"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返回条数，若传入超过系统限制，返回系统限制数量"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/messages",targetCode:`curl -X GET '${i.appDetail.api_base_url}/messages?user=abc-123&conversation_id={conversation_id}'
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
"limit": 20,
"has_more": false,
"data": [
    {
        "id": "a076a87f-31e5-48dc-b452-0061adbbc922",
        "conversation_id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
        "inputs": {
            "name": "gofy"
        },
        "query": "iphone 13 pro",
        "answer": "The iPhone 13 Pro, released on September 24, 2021, features a 6.1-inch display with a resolution of 1170 x 2532. It is equipped with a Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard) processor, 6 GB of RAM, and offers storage options of 128 GB, 256 GB, 512 GB, and 1 TB. The camera is 12 MP, the battery capacity is 3095 mAh, and it runs on iOS 15.",
        "message_files": [],
        "feedback": null,
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ],
        "created_at": 1705569239
    }
  ]
}
`})})}),e.jsx(n.h3,{children:"Response Example(智能助手)"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
"limit": 20,
"has_more": false,
"data": [
    {
        "id": "d35e006c-7c4d-458f-9142-be4930abdf94",
        "conversation_id": "957c068b-f258-4f89-ba10-6e8a0361c457",
        "inputs": {},
        "query": "draw a cat",
        "answer": "I have generated an image of a cat for you. Please check your messages to view the image.",
        "message_files": [
            {
                "id": "976990d2-5294-47e6-8f14-7356ba9d2d76",
                "type": "image",
                "url": "http://127.0.0.1:5001/files/tools/976990d2-5294-47e6-8f14-7356ba9d2d76.png?timestamp=1705988524&nonce=55df3f9f7311a9acd91bf074cd524092&sign=z43nMSO1L2HBvoqADLkRxr7Biz0fkjeDstnJiCK1zh8=",
                "belongs_to": "assistant"
            }
        ],
        "feedback": null,
        "retriever_resources": [],
        "created_at": 1705988187
    }
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations",method:"GET",title:"获取会话列表",name:"#conversations"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"获取当前用户的会话列表，默认返回最近的 20 条。"}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"（选填）当前页最后面一条记录的 ID，默认 null"})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"（选填）一次请求返回多少条记录，默认 20 条，最大 100 条，最小 1 条。"})},"limit"),e.jsxs(d,{name:"sort_by",type:"string",children:[e.jsx(n.p,{children:"（选填）排序字段，默认 -updated_at(按更新时间倒序排列)"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"可选值：created_at, -created_at, updated_at, -updated_at"}),`
`,e.jsx(n.li,{children:"字段前面的符号代表顺序或倒序，-代表倒序"}),`
`]})]},"sort_by")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 会话列表",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"}),"  (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"}),"  (string) 会话名称，默认由大语言模型生成。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 用户输入参数。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 会话状态"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) 开场白"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 创建时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) 更新时间"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返回条数，若传入超过系统限制，返回系统限制数量"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/conversations",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations?user=abc-123&last_id=&limit=20' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
      "id": "10799fb8-64f7-4296-bbf7-b42bfbe0ae54",
      "name": "New chat",
      "inputs": {
          "book": "book",
          "myName": "Lucy"
      },
      "status": "normal",
      "created_at": 1679667915,
      "updated_at": 1679667915
    },
    {
      "id": "hSIhXBhNe8X1d8Et"
      // ...
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id",method:"DELETE",title:"删除会话",name:"#delete"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"删除会话。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"})," (string) 固定返回 success"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"DELETE",label:"/conversations/:conversation_id",targetCode:`curl -X DELETE '${i.appDetail.api_base_url}/conversations/{conversation_id}' \\
--header 'Content-Type: application/json' \\
--header 'Accept: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data '{
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/name",method:"POST",title:"会话重命名",name:"#rename"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"对会话进行重命名，会话名称用于显示在支持多会话的客户端上。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"name",type:"string",children:e.jsxs(n.p,{children:["（选填）名称，若 ",e.jsx(n.code,{children:"auto_generate"})," 为 ",e.jsx(n.code,{children:"true"})," 时，该参数可不传。"]})},"name"),e.jsx(d,{name:"auto_generate",type:"bool",children:e.jsx(n.p,{children:"（选填）自动生成标题，默认 false。"})},"auto_generate"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"}),"  (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"}),"  (string) 会话名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 用户输入参数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 会话状态"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) 开场白"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 创建时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) 更新时间"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/conversations/:conversation_id/name",targetCode:`curl -X POST '${i.appDetail.api_base_url}/conversations/{conversation_id}/name' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "name": "",
  "auto_generate": true,
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "34d511d5-56de-4f16-a997-57b379508443",
  "name": "hello",
  "inputs": {},
  "status": "normal",
  "introduction": "",
  "created_at": 1732731141,
  "updated_at": 1732734510
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables",method:"GET",title:"获取对话变量",name:"#conversation-variables"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"从特定对话中检索变量。此端点对于提取对话过程中捕获的结构化数据非常有用。"}),e.jsx(n.h3,{children:"路径参数"}),e.jsx(t,{children:e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"要从中检索变量的对话ID。"})},"conversation_id")}),e.jsx(n.h3,{children:"查询参数"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识符，由开发人员定义的规则，在应用程序内必须唯一。"})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"（选填）当前页最后面一条记录的 ID，默认 null"})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"（选填）一次请求返回多少条记录，默认 20 条，最大 100 条，最小 1 条。"})},"limit")]}),e.jsx(n.h3,{children:"响应"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 每页项目数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) 是否有更多项目"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 变量列表",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 变量 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 变量名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) 变量类型（字符串、数字、布尔等）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (string) 变量值"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 变量描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) 最后更新时间戳"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"错误"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", 对话不存在"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/conversations/:conversation_id/variables",debug:"true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"带变量名过滤的请求",language:"bash",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123&variable_name=customer_name' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 100,
  "has_more": false,
  "data": [
    {
      "id": "variable-uuid-1",
      "name": "customer_name",
      "value_type": "string",
      "value": "John Doe",
      "description": "客户名称（从对话中提取）",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    },
    {
      "id": "variable-uuid-2",
      "name": "order_details",
      "value_type": "json",
      "value": "{\\"product\\":\\"Widget\\",\\"quantity\\":5,\\"price\\":19.99}",
      "description": "客户的订单详情",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables/:variable_id",method:"PUT",title:"更新对话变量",name:"#update-conversation-variable"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"更新特定对话变量的值。此端点允许您修改在对话过程中捕获的变量值，同时保留其名称、类型和描述。"}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"包含要更新变量的对话ID。"})},"conversation_id"),e.jsx(d,{name:"variable_id",type:"string",children:e.jsx(n.p,{children:"要更新的变量ID。"})},"variable_id")]}),e.jsx(n.h3,{children:"请求体"}),e.jsxs(t,{children:[e.jsx(d,{name:"value",type:"any",children:e.jsx(n.p,{children:"变量的新值。必须匹配变量的预期类型（字符串、数字、对象等）。"})},"value"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识符，由开发人员定义的规则，在应用程序内必须唯一。"})},"user")]}),e.jsx(n.h3,{children:"响应"}),e.jsx(n.p,{children:"返回包含以下内容的更新变量对象："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 变量ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 变量名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) 变量类型（字符串、数字、对象等）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (any) 更新后的变量值"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 变量描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) 最后更新时间戳"]}),`
`]}),e.jsx(n.h3,{children:"错误"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"Type mismatch: variable expects {expected_type}, but got {actual_type} type"}),", 值类型与变量的预期类型不匹配"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", 对话不存在"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_variable_not_exists"}),", 变量不存在"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/conversations/:conversation_id/variables/:variable_id",targetCode:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": "Updated Value",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"使用不同值类型更新",targetCode:[{title:"字符串值",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": "新的字符串值",
  "user": "abc-123"
}'`},{title:"数字值",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": 42,
  "user": "abc-123"
}'`},{title:"对象值",code:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
  "value": {"product": "Widget", "quantity": 10, "price": 29.99},
  "user": "abc-123"
}'`}]}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "variable-uuid-1",
  "name": "customer_name",
  "value_type": "string",
  "value": "Updated Value",
  "description": "客户名称（从对话中提取）",
  "created_at": 1650000000000,
  "updated_at": 1650000001000
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/audio-to-text",method:"POST",title:"语音转文字",name:"#audio-to-text"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["该接口需使用 ",e.jsx(n.code,{children:"multipart/form-data"})," 进行请求。"]}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsxs(n.p,{children:[`语音文件。
支持格式：`,e.jsx(n.code,{children:"['mp3', 'mp4', 'mpeg', 'mpga', 'm4a', 'wav', 'webm']"}),`
文件大小限制：15MB`]})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) 输出文字"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/audio-to-text",targetCode:`curl -X POST '${i.appDetail.api_base_url}/audio-to-text' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=audio/[mp3|mp4|mpeg|mpga|m4a|wav|webm]'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "text": "hello"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"文字转语音",name:"#text-to-audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"文字转语音。"}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"Gofy 生成的文本消息，那么直接传递生成的message-id 即可，后台会通过 message_id 查找相应的内容直接合成语音信息。如果同时传 message_id 和 text，优先使用 message_id。"})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"语音生成内容。如果没有传 message-id的话，则会使用这个字段的内容"})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/text-to-audio",targetCode:`curl -o text-to-audio.mp3 -X POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290",
  "text": "Hello Gofy",
  "user": "abc-123",
}'`}),e.jsx(s,{title:"headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"获取应用基本信息",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取应用的基本信息"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 应用名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 应用描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) 应用标签"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "advanced-chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"获取应用参数",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于进入页面一开始，获取功能开关、输入参数名称、类型及默认值等使用。"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) 开场白"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) 开场推荐问题列表"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) 启用回答后给出推荐问题。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) 语音转文本",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text_to_speech"})," (object) 文本转语音",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"voice"})," (string) 语音类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"language"})," (string) 语言"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"autoPlay"})," (string) 自动播放",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"}),"  开启"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"disabled"})," 关闭"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) 引用和归属",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) 标记回复",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) 用户输入表单配置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) 文本输入控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) 段落文本输入控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) 下拉控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) 选项值"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) 文件上传配置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) 文档设置
当前仅支持文档类型：`,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 文档数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) 图片设置
当前仅支持图片类型：`,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 图片数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) 音频设置
当前仅支持音频类型：`,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 音频数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) 视频设置
当前仅支持视频类型：`,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 视频数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) 自定义设置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 自定义数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) 系统参数",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) Document upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) Image file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) Audio file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) Video file upload size limit (MB)"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:`curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "introduction": "nice to meet you",
  "user_input_form": [
    {
      "text-input": {
        "label": "a",
        "variable": "a",
        "required": true,
        "max_length": 48,
        "default": ""
      }
    },
    {
      // ...
    }
  ],
  "file_upload": {
    "image": {
      "enabled": true,
      "number_limits": 3,
      "transfer_methods": [
        "remote_url",
        "local_file"
      ]
    }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/meta",method:"GET",title:"获取应用Meta信息",name:"#meta"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取工具 icon"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_icons"}),"(object[string]) 工具图标",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"工具名称"})," (string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (object|string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["(object) 图标",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"background"})," (string) hex 格式的背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"content"}),"(string) emoji"]}),`
`]}),`
`]}),`
`,e.jsx(n.li,{children:"(string) 图标 URL"}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/meta",targetCode:`curl -X GET '${i.appDetail.api_base_url}/meta' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "tool_icons": {
      "dalle2": "https://cloud.gofy.ai/console/api/workspaces/current/tool-provider/builtin/dalle/icon",
      "api_tool": {
          "background": "#252525",
          "content": "\\ud83d\\ude01"
      }
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"获取应用 WebApp 设置",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取应用的 WebApp 设置"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp 名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) 聊天颜色主题，hex 格式"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) 聊天颜色主题是否反转"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) 图标类型，",e.jsx(n.code,{children:"emoji"}),"-表情，",e.jsx(n.code,{children:"image"}),"-图片"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) 图标，如果是 ",e.jsx(n.code,{children:"emoji"})," 类型，则是 emoji 表情符号，如果是 ",e.jsx(n.code,{children:"image"})," 类型，则是图片 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) hex 格式的背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) 图标 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) 版权信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) 隐私政策链接"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) 自定义免责声明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) 默认语言"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) 是否显示工作流详情"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) 是否使用 WebApp 图标替换聊天中的 🤖"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations",method:"GET",title:"获取标注列表",name:"#annotation_list"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"页码"})},"page"),e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"每页数量"})},"limit")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/apps/annotations",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/apps/annotations?page=1&limit=20' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "data": [
    {
      "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
      "question": "What is your name?",
      "answer": "I am Gofy.",
      "hit_count": 0,
      "created_at": 1735625869
    }
  ],
  "has_more": false,
  "limit": 20,
  "total": 1,
  "page": 1
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations",method:"POST",title:"创建标注",name:"#create_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"question",type:"string",children:e.jsx(n.p,{children:"问题"})},"question"),e.jsx(d,{name:"answer",type:"string",children:e.jsx(n.p,{children:"答案内容"})},"answer")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/apps/annotations",targetCode:`curl --location --request POST '${i.appDetail.api_base_url}/apps/annotations' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "question": "What is your name?",
  "answer": "I am Gofy."
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
  "question": "What is your name?",
  "answer": "I am Gofy.",
  "hit_count": 0,
  "created_at": 1735625869
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations/{annotation_id}",method:"PUT",title:"更新标注",name:"#update_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"annotation_id",type:"string",children:e.jsx(n.p,{children:"标注 ID"})},"annotation_id"),e.jsx(d,{name:"question",type:"string",children:e.jsx(n.p,{children:"问题"})},"question"),e.jsx(d,{name:"answer",type:"string",children:e.jsx(n.p,{children:"答案内容"})},"answer")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/apps/annotations/{annotation_id}",targetCode:`curl --location --request PUT '${i.appDetail.api_base_url}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "question": "What is your name?",
  "answer": "I am Gofy."
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
  "question": "What is your name?",
  "answer": "I am Gofy.",
  "hit_count": 0,
  "created_at": 1735625869
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations/{annotation_id}",method:"DELETE",title:"删除标注",name:"#delete_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"annotation_id",type:"string",children:e.jsx(n.p,{children:"标注 ID"})},"annotation_id")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"DELETE",label:"/apps/annotations/{annotation_id}",targetCode:`curl --location --request DELETE '${i.appDetail.api_base_url}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotation-reply/{action}",method:"POST",title:"标注回复初始设置",name:"#initial_annotation_reply_settings"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"action",type:"string",children:e.jsx(n.p,{children:"动作，只能是 'enable' 或 'disable'"})},"action"),e.jsx(d,{name:"embedding_provider_name",type:"string",children:e.jsx(n.p,{children:"指定的嵌入模型提供商，必须先在系统内设定好接入的模型，对应的是 provider 字段"})},"embedding_provider_name"),e.jsx(d,{name:"embedding_model_name",type:"string",children:e.jsx(n.p,{children:"指定的嵌入模型，对应的是 model 字段"})},"embedding_model_name"),e.jsx(d,{name:"score_threshold",type:"number",children:e.jsx(n.p,{children:"相似度阈值，当相似度大于该阈值时，系统会自动回复，否则不回复"})},"score_threshold")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.p,{children:"嵌入模型的提供商和模型名称可以通过以下接口获取：v1/workspaces/current/models/model-types/text-embedding，具体见：通过 API 维护知识库。使用的 Authorization 是 Dataset 的 API Token。"}),e.jsx(s,{title:"Request",tag:"POST",label:"/apps/annotation-reply/{action}",targetCode:`curl --location --request POST '${i.appDetail.api_base_url}/apps/annotation-reply/{action}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "score_threshold": 0.9,
  "embedding_provider_name": "zhipu",
  "embedding_model_name": "embedding_3"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "job_id": "b15c8f68-1cf4-4877-bf21-ed7cf2011802",
  "job_status": "waiting"
}
`})})}),e.jsx(n.p,{children:"该接口是异步执行，所以会返回一个job_id，通过查询job状态接口可以获取到最终的执行结果。"})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotation-reply/{action}/status/{job_id}",method:"GET",title:"查询标注回复初始设置任务状态",name:"#initial_annotation_reply_settings_task_status"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"action",type:"string",children:e.jsx(n.p,{children:"动作，只能是 'enable' 或 'disable'，并且必须和标注回复初始设置接口的动作一致"})},"action"),e.jsx(d,{name:"job_id",type:"string",children:e.jsx(n.p,{children:"任务 ID，从标注回复初始设置接口返回的 job_id"})},"job_id")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/apps/annotations",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/apps/annotation-reply/{action}/status/{job_id}' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "job_id": "b15c8f68-1cf4-4877-bf21-ed7cf2011802",
  "job_status": "waiting",
  "error_msg": ""
}
`})})})]})]})]})}function Kn(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(_e,{...i})}):_e(i)}function me(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"Chat App API"}),`
`,e.jsx(n.p,{children:"Chat applications support session persistence, allowing previous chat history to be used as context for responses. This can be applicable for chatbot, customer service AI, etc."}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"Base URL"}),e.jsx(s,{title:"Code",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"Authentication"}),e.jsxs(n.p,{children:["The Service API uses ",e.jsx(n.code,{children:"API-Key"}),` authentication.
`,e.jsx("i",{children:e.jsx(n.strong,{children:"Strongly recommend storing your API Key on the server-side, not shared or stored on the client-side, to avoid possible API-Key leakage that can lead to serious consequences."})})]}),e.jsxs(n.p,{children:["For all API requests, include your API Key in the ",e.jsx(n.code,{children:"Authorization"}),"HTTP Header, as shown below:"]}),e.jsx(s,{title:"Code",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages",method:"POST",title:"Send Chat Message",name:"#Send-Chat-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Send a request to the chat application."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"query",type:"string",children:e.jsx(n.p,{children:"User Input/Question content"})},"query"),e.jsx(d,{name:"inputs",type:"object",children:e.jsxs(n.p,{children:[`Allows the entry of various variable values defined by the App.
The `,e.jsx(n.code,{children:"inputs"})," parameter contains multiple key/value pairs, with each key corresponding to a specific variable and each value being the specific value for that variable. Default ",e.jsx(n.code,{children:"{}"})]})},"inputs"),e.jsxs(d,{name:"response_mode",type:"string",children:[e.jsx(n.p,{children:"The mode of response return, supporting:"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," Streaming mode (recommended), implements a typewriter-like output through SSE (",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` Blocking mode, returns result after execution is complete. (Requests may be interrupted if the process is long)
Due to Cloudflare restrictions, the request will be interrupted without a return after 100 seconds.
`,e.jsx("i",{children:"Note: blocking mode is not supported in Agent Assistant mode"})]}),`
`]})]},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`})},"user"),e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"Conversation ID, to continue the conversation based on previous chat records, it is necessary to pass the previous message's conversation_id."})},"conversation_id"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"File list, suitable for inputting files combined with text understanding and answering questions, available only when the model supports Vision/Video capability."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) Supported type:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," Supported types include: 'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," Supported types include: 'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," Supported types include: 'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," Supported types include: 'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," Supported types include: other file types"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) Transfer method:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": File URL."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": Upload file."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," File URL. (Only when transfer method is ",e.jsx(n.code,{children:"remote_url"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," Upload file ID. (Only when transfer method is ",e.jsx(n.code,{children:"local_file"}),")."]}),`
`]})]},"files"),e.jsx(d,{name:"auto_generate_name",type:"bool",children:e.jsxs(n.p,{children:["Auto-generate title, default is ",e.jsx(n.code,{children:"true"}),`.
If set to `,e.jsx(n.code,{children:"false"}),", can achieve async title generation by calling the conversation rename API and setting ",e.jsx(n.code,{children:"auto_generate"})," to ",e.jsx(n.code,{children:"true"}),"."]})},"auto_generate_name"),e.jsx(d,{name:"workflow_id",type:"string",children:e.jsxs(n.p,{children:["(Optional) Workflow ID to specify a specific version, if not provided, uses the default published version.",e.jsx("br",{}),`
How to obtain: In the version history interface, click the copy icon on the right side of each version entry to copy the complete workflow ID.`]})},"workflow_id"),e.jsxs(d,{name:"trace_id",type:"string",children:[e.jsxs(n.p,{children:["(Optional) Trace ID. Used for integration with existing business trace components to achieve end-to-end distributed tracing. If not provided, the system will automatically generate a trace_id. Supports the following three ways to pass, in order of priority:",e.jsx("br",{})]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["Header: via HTTP Header ",e.jsx("code",{children:"X-Trace-Id"}),", highest priority.",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["Query parameter: via URL query parameter ",e.jsx("code",{children:"trace_id"}),".",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["Request Body: via request body field ",e.jsx("code",{children:"trace_id"})," (i.e., this field).",e.jsx("br",{})]}),`
`]})]},"trace_id")]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:`When response_mode is blocking, return a CompletionResponse object.
When response_mode is streaming, return a ChunkCompletionResponse stream.`}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["Returns the complete App result, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"application/json"}),"."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) Event type, always ",e.jsx(n.code,{children:"message"})," in blocking mode."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID, same as ",e.jsx(n.code,{children:"message_id"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) App mode, fixed as ",e.jsx(n.code,{children:"chat"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) Complete response content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) Metadata",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) Model usage information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) Citation and Attribution List"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Message creation timestamp, e.g., 1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["Returns the stream chunks outputted by the App, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"text/event-stream"}),`.
Each streaming chunk starts with `,e.jsx(n.code,{children:"data:"}),", separated by two newline characters ",e.jsx(n.code,{children:"\\n\\n"}),", as shown below:"]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsxs(n.p,{children:["The structure of the streaming chunks varies depending on the ",e.jsx(n.code,{children:"event"}),":"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLM returns text chunk event, i.e., the complete text is output in a chunked fashion.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLM returned text chunk content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: agent_message"})," LLM returns text chunk event, i.e., with Agent Assistant enabled, the complete text is output in a chunked fashion (Only supported in Agent mode)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLM returned text chunk content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS audio stream event, that is, speech synthesis output. The content is an audio block in Mp3 format, encoded as a base64 string. When playing, simply decode the base64 and feed it into the player. (This message is available only when auto-play is enabled)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the stop response interface below"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The audio after speech synthesis, encoded in base64 text content, when playing, simply decode the base64 and feed it into the player"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g.: 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS audio stream end event, receiving this event indicates the end of the audio stream.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the stop response interface below"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The end event has no audio, so this is an empty string"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g.: 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: agent_thought"})," thought of Agent, contains the thought of LLM, input and output of tool calls (Only supported in Agent mode)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Agent thought ID, every iteration has a unique agent thought ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string)  Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"position"})," (int) Position of current agent thought, each message may have multiple thoughts in order."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"thought"})," (string) What LLM is thinking about"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"observation"})," (string) Response from tool calls"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool"})," (string) A list of tools represents which tools are called，split by ;"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_input"})," (string) Input of tools in JSON format. Like: ",e.jsx(n.code,{children:'{"dalle3": {"prompt": "a cute cat"}}'}),"."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[string])  Refer to message_file event",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"})," (string) File ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_file"})," Message file event, a new file has created by tool",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) File unique ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"}),' (string) File type，only allow "image" currently']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) Belongs to, it will only be an 'assistant' here"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) Remote url of file"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"}),"  (string) Conversation ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," Message end event, receiving this event means streaming has ended.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) Metadata",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) Model usage information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) Citation and Attribution List"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` Message content replacement event.
When output content moderation is enabled, if the content is flagged, then the message content will be replaced with a preset reply through this event.`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) Replacement content (directly replaces all LLM reply text)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
Exceptions that occur during the streaming process will be output in the form of stream events, and reception of an error event will end the stream.`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTP status code"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) Error code"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) Error message"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," Ping event every 10 seconds to keep the connection alive."]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404, Conversation does not exists"}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", abnormal parameter input"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", App configuration unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", no available model credential configuration"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", model invocation quota insufficient"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", current model unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_not_found"}),", specified workflow version not found"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"draft_workflow_error"}),", cannot use draft workflow version"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_id_format_error"}),", invalid workflow_id format, expected UUID format"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"completion_request_error"}),", text generation failed"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/chat-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "query": "What are the specs of the iPhone 13 Pro Max?",
  "response_mode": "streaming",
  "conversation_id": "",
  "user": "abc-123",
  "files": [
      {
        "type": "image",
      "transfer_method": "remote_url",
      "url": "https://cloud.gofy.ai/logo/logo-site.png"
    }
  ]
}'`}),e.jsx(n.h3,{children:"Blocking Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "event": "message",
    "task_id": "c3800678-a077-43df-a102-53f23ed20b88",
    "id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "message_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
    "mode": "chat",
    "answer": "iPhone 13 Pro Max specs are listed here:...",
    "metadata": {
        "usage": {
            "prompt_tokens": 1033,
            "prompt_unit_price": "0.001",
            "prompt_price_unit": "0.001",
            "prompt_price": "0.0010330",
            "completion_tokens": 128,
            "completion_unit_price": "0.002",
            "completion_price_unit": "0.001",
            "completion_price": "0.0002560",
            "total_tokens": 1161,
            "total_price": "0.0012890",
            "currency": "USD",
            "latency": 0.7682376249867957
        },
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ]
    },
    "created_at": 1705407629
}
`})})}),e.jsx(n.h3,{children:"Streaming Mode ( Basic Assistant )"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " I", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": "'m", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " glad", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " to", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " meet", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " you", "created_at": 1679586595}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})}),e.jsx(n.h3,{children:"Response Example(Agent Assistant)"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " I", "created_at": 1679586595}
data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": "'m", "created_at": 1679586595}
data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " glad", "created_at": 1679586595}
data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " to", "created_at": 1679586595}
data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " meet", "created_at": 1679586595}
data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " you", "created_at": 1679586595}
data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"File Upload",name:"#file-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:`Upload a file (currently only images are supported) for use when sending messages, enabling multimodal understanding of images and text.
Supports png, jpg, jpeg, webp, gif formats.
Uploaded files are for use by the current end-user only.`}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["This interface requires a ",e.jsx(n.code,{children:"multipart/form-data"})," request."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file"}),` (File) Required
The file to be uploaded.`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, defined by the developer's rules, must be unique within the application. The Service API does not share conversations created by the WebApp.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"After a successful upload, the server will return the file's ID and related information."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) File name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) File size (bytes)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) File extension"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) File mime-type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) End-user ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"no_file_uploaded"}),", a file must be provided"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"too_many_files"}),", currently only one file is accepted"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_preview"}),", the file does not support preview"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_estimate"}),", the file does not support estimation"]}),`
`,e.jsxs(n.li,{children:["413, ",e.jsx(n.code,{children:"file_too_large"}),", the file is too large"]}),`
`,e.jsxs(n.li,{children:["415, ",e.jsx(n.code,{children:"unsupported_file_type"}),", unsupported extension, currently only document files are accepted"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_connection_failed"}),", unable to connect to S3 service"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_permission_denied"}),", no permission to upload files to S3"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_file_too_large"}),", file exceeds S3 size limit"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"Get End User",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve an end user by ID."}),e.jsxs(n.p,{children:["This is useful when other APIs return an end-user ID (e.g. ",e.jsx(n.code,{children:"created_by"})," from File Upload)."]}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) Required
End user ID.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns an EndUser object."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) Tenant ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) App ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) End user type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) External user ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) Whether anonymous"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) Session ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 datetime"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 datetime"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"end_user_not_found"}),", end user not found"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"File Preview",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Preview or download uploaded files. This endpoint allows you to access files that have been previously uploaded via the File Upload API."}),e.jsx("i",{children:"Files can only be accessed if they belong to messages within the requesting application."}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) Required
The unique identifier of the file to preview, obtained from the File Upload API response.`]}),`
`]}),e.jsx(n.h3,{children:"Query Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) Optional
Whether to force download the file as an attachment. Default is `,e.jsx(n.code,{children:"false"})," (preview in browser)."]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns the file content with appropriate headers for browser display or download."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," Set based on file mime type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," File size in bytes (if available)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"}),' Set to "attachment" if ',e.jsx(n.code,{children:"as_attachment=true"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," Caching headers for performance"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' Set to "bytes" for audio/video files']}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", abnormal parameter input"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", file access denied or file does not belong to current application"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", file not found or has been deleted"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Download as Attachment"}),e.jsx(s,{title:"Download Request",tag:"GET",label:"/files/:file_id/preview?as_attachment=true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"Response Headers Example"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"Download Response Headers"}),e.jsx(s,{title:"Download Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages/:task_id/stop",method:"POST",title:"Stop Generate",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Only supported in streaming mode."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, can be obtained from the streaming chunk return"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, used to define the identity of the end-user, must be consistent with the user passed in the message sending interface. The Service API does not share conversations created by the WebApp.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/chat-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{"user": "abc-123"}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"Message Feedback",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"End-users can provide feedback messages, facilitating application developers to optimize expected outputs."}),e.jsx(n.h3,{children:"Path"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"Message ID"})},"message_id")}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsxs(n.p,{children:["Upvote as ",e.jsx(n.code,{children:"like"}),", downvote as ",e.jsx(n.code,{children:"dislike"}),", revoke upvote as ",e.jsx(n.code,{children:"null"})]})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"User identifier, defined by the developer's rules, must be unique within the application."})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"The specific content of message feedback."})},"content")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"Get feedbacks of application",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Get application's feedbacks."}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（optional）pagination，default：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（optional） records per page default：20"})},"limit")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (List) return apps feedback list."]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`  {
      "data": [
          {
              "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
              "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
              "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
              "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
              "rating": "like",
              "content": "message feedback information-3",
              "from_source": "user",
              "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
              "from_account_id": null,
              "created_at": "2025-04-24T09:24:38",
              "updated_at": "2025-04-24T09:24:38"
          }
      ]
  }
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/{message_id}/suggested",method:"GET",title:"Next Suggested Questions",name:"#suggested"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Get next questions suggestions for the current message"}),e.jsx(n.h3,{children:"Path Params"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"Message ID"})},"message_id")}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`})},"user")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/messages/{message_id}/suggested",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/messages/{message_id}/suggested?user=abc-123& \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success",
  "data": [
        "a",
        "b",
        "c"
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages",method:"GET",title:"Get Conversation History Messages",name:"#messages"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:["Returns historical chat records in a scrolling load format, with the first page returning the latest ",e.jsx(n.code,{children:"{limit}"})," messages, i.e., in reverse order."]}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"Conversation ID"})},"conversation_id"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`})},"user"),e.jsx(d,{name:"first_id",type:"string",children:e.jsx(n.p,{children:"The ID of the first chat record on the current page, default is null."})},"first_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"How many chat history messages to return in one request, default is 20."})},"limit")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) Message list",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) User input parameters."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"})," (string) User input / question content."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[object]) Message files",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) File type, image for images"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) File preview URL, use the File Preview API (",e.jsx(n.code,{children:"/files/{file_id}/preview"}),") to access the file"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) belongs to，user or assistant"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"agent_thoughts"})," (array[object]) Agent thought（Empty if it's a Basic Assistant）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Agent thought ID, every iteration has a unique agent thought ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"position"})," (int) Position of current agent thought, each message may have multiple thoughts in order."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"thought"})," (string) What LLM is thinking about"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"observation"})," (string) Response from tool calls"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool"})," (string) A list of tools represents which tools are called，split by ;"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_input"})," (string) Input of tools in JSON format. Like: ",e.jsx(n.code,{children:'{"dalle3": {"prompt": "a cute cat"}}'}),"."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[string])  Refer to message_file event",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"})," (string) File ID"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) Response message content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"feedback"})," (object) Feedback information",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"rating"})," (string) Upvote as ",e.jsx(n.code,{children:"like"})," / Downvote as ",e.jsx(n.code,{children:"dislike"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) Citation and Attribution List"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) Whether there is a next page"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) Number of returned items, if input exceeds system limit, returns system limit amount"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/messages",targetCode:`curl -X GET '${i.appDetail.api_base_url}/messages?user=abc-123&conversation_id='\\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example (Basic Assistant)"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
        "id": "a076a87f-31e5-48dc-b452-0061adbbc922",
        "conversation_id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
        "inputs": {
            "name": "gofy"
        },
        "query": "iphone 13 pro",
        "answer": "The iPhone 13 Pro, released on September 24, 2021, features a 6.1-inch display with a resolution of 1170 x 2532. It is equipped with a Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard) processor, 6 GB of RAM, and offers storage options of 128 GB, 256 GB, 512 GB, and 1 TB. The camera is 12 MP, the battery capacity is 3095 mAh, and it runs on iOS 15.",
        "message_files": [],
        "feedback": null,
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ],
        "agent_thoughts": [],
        "created_at": 1705569239,
    }
  ]
}
`})})}),e.jsx(n.h3,{children:"Response Example (Agent Assistant)"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "limit": 20,
    "has_more": false,
    "data": [
        {
            "id": "d35e006c-7c4d-458f-9142-be4930abdf94",
            "conversation_id": "957c068b-f258-4f89-ba10-6e8a0361c457",
            "inputs": {},
            "query": "draw a cat",
            "answer": "I have generated an image of a cat for you. Please check your messages to view the image.",
            "message_files": [
                {
                    "id": "976990d2-5294-47e6-8f14-7356ba9d2d76",
                    "type": "image",
                    "url": "http://127.0.0.1:5001/files/tools/976990d2-5294-47e6-8f14-7356ba9d2d76.png?timestamp=1705988524&nonce=55df3f9f7311a9acd91bf074cd524092&sign=z43nMSO1L2HBvoqADLkRxr7Biz0fkjeDstnJiCK1zh8=",
                    "belongs_to": "assistant"
                }
            ],
            "feedback": null,
            "retriever_resources": [],
            "created_at": 1705988187,
            "agent_thoughts": [
                {
                    "id": "592c84cf-07ee-441c-9dcc-ffc66c033469",
                    "chain_id": null,
                    "message_id": "d35e006c-7c4d-458f-9142-be4930abdf94",
                    "position": 1,
                    "thought": "",
                    "tool": "dalle2",
                    "tool_input": "{\\"dalle2\\": {\\"prompt\\": \\"cat\\"}}",
                    "created_at": 1705988186,
                    "observation": "image has been created and sent to user already, you should tell user to check it now.",
                    "files": [
                        "976990d2-5294-47e6-8f14-7356ba9d2d76"
                    ]
                },
                {
                    "id": "73ead60d-2370-4780-b5ed-532d2762b0e5",
                    "chain_id": null,
                    "message_id": "d35e006c-7c4d-458f-9142-be4930abdf94",
                    "position": 2,
                    "thought": "I have generated an image of a cat for you. Please check your messages to view the image.",
                    "tool": "",
                    "tool_input": "",
                    "created_at": 1705988199,
                    "observation": "",
                    "files": []
                }
            ]
        }
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations",method:"GET",title:"Get Conversations",name:"#conversations"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve the conversation list for the current user, defaulting to the most recent 20 entries."}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"(Optional) The ID of the last record on the current page, default is null."})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"(Optional) How many records to return in one request, default is the most recent 20 entries. Maximum 100, minimum 1."})},"limit"),e.jsxs(d,{name:"sort_by",type:"string",children:[e.jsx(n.p,{children:"(Optional) Sorting Field, Default: -updated_at (sorted in descending order by update time)"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"Available Values: created_at, -created_at, updated_at, -updated_at"}),`
`,e.jsx(n.li,{children:'The symbol before the field represents the order or reverse, "-" represents reverse order.'}),`
`]})]},"sort_by")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) List of conversations",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Conversation name, by default, is a snippet of the first question asked by the user in the conversation."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) User input parameters."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) Conversation status"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) Introduction"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) Update timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) Number of entries returned, if input exceeds system limit, system limit number is returned"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/conversations",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations?user=abc-123&last_id=&limit=20' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
      "id": "10799fb8-64f7-4296-bbf7-b42bfbe0ae54",
      "name": "New chat",
      "inputs": {
          "book": "book",
          "myName": "Lucy"
      },
      "status": "normal",
      "created_at": 1679667915,
      "updated_at": 1679667915
    },
    {
      "id": "hSIhXBhNe8X1d8Et"
      // ...
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id",method:"DELETE",title:"Delete Conversation",name:"#delete"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Delete a conversation."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the application."})},"user")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"DELETE",label:"/conversations/:conversation_id",targetCode:`curl -X DELETE '${i.appDetail.api_base_url}/conversations/:conversation_id' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/name",method:"POST",title:"Conversation Rename",name:"#rename"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Request Body"}),e.jsx(n.p,{children:"Rename the session, the session name is used for display on clients that support multiple sessions."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) Conversation ID"]}),`
`]}),e.jsxs(t,{children:[e.jsx(d,{name:"name",type:"string",children:e.jsxs(n.p,{children:["(Optional) The name of the conversation. This parameter can be omitted if ",e.jsx(n.code,{children:"auto_generate"})," is set to ",e.jsx(n.code,{children:"true"}),"."]})},"name"),e.jsx(d,{name:"auto_generate",type:"bool",children:e.jsxs(n.p,{children:["(Optional) Automatically generate the title, default is ",e.jsx(n.code,{children:"false"})]})},"auto_generate"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the application."})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Conversation ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Conversation name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) User input parameters"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) Conversation status"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) Introduction"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) Update timestamp, e.g., 1705395332"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/conversations/:conversation_id/name",targetCode:`curl -X POST '${i.appDetail.api_base_url}/conversations/:conversation_id/name' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "name": "",
  "auto_generate": true,
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
    "name": "Chat vs AI",
    "inputs": {},
    "status": "normal",
    "introduction": "",
    "created_at": 1705569238,
    "updated_at": 1705569238
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables",method:"GET",title:"Get Conversation Variables",name:"#conversation-variables"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve variables from a specific conversation. This endpoint is useful for extracting structured data that was captured during the conversation."}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsx(t,{children:e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"The ID of the conversation to retrieve variables from."})},"conversation_id")}),e.jsx(n.h3,{children:"Query Parameters"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the application"})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"(Optional) The ID of the last record on the current page, default is null."})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"(Optional) How many records to return in one request, default is the most recent 20 entries. Maximum 100, minimum 1."})},"limit")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) Number of items per page"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) Whether there is a next page"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) List of variables",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Variable name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) Variable type (string, number, object, etc.)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (string) Variable value"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) Variable description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) Last update timestamp"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", Conversation not found"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/conversations/:conversation_id/variables",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Request with variable name filter",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X GET '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123&variable_name=customer_name' \\
--header 'Authorization: Bearer {api_key}'
`})})}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 100,
  "has_more": false,
  "data": [
    {
      "id": "variable-uuid-1",
      "name": "customer_name",
      "value_type": "string",
      "value": "John Doe",
      "description": "Customer name extracted from the conversation",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    },
    {
      "id": "variable-uuid-2",
      "name": "order_details",
      "value_type": "json",
      "value": "{\\"product\\":\\"Widget\\",\\"quantity\\":5,\\"price\\":19.99}",
      "description": "Order details from the customer",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables/:variable_id",method:"PUT",title:"Update Conversation Variable",name:"#update-conversation-variable"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Update the value of a specific conversation variable. This endpoint allows you to modify the value of a variable that was captured during the conversation while preserving its name, type, and description."}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"The ID of the conversation containing the variable to update."})},"conversation_id"),e.jsx(d,{name:"variable_id",type:"string",children:e.jsx(n.p,{children:"The ID of the variable to update."})},"variable_id")]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"value",type:"any",children:e.jsx(n.p,{children:"The new value for the variable. Must match the variable's expected type (string, number, object, etc.)."})},"value"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the application."})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns the updated variable object with:"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Variable name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) Variable type (string, number, object, etc.)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (any) Updated variable value"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) Variable description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) Last update timestamp"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"Type mismatch: variable expects {expected_type}, but got {actual_type} type"}),", Value type doesn't match variable's expected type"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", Conversation not found"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_variable_not_exists"}),", Variable not found"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/conversations/:conversation_id/variables/:variable_id",targetCode:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "value": "Updated Value",
  "user": "abc-123"
}'`}),e.jsxs(s,{title:"Update with different value types",children:[e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": "New string value",
    "user": "abc-123"
}'
`})}),e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": 42,
    "user": "abc-123"
}'
`})}),e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": {"product": "Widget", "quantity": 10, "price": 29.99},
    "user": "abc-123"
}'
`})})]}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "variable-uuid-1",
  "name": "customer_name",
  "value_type": "string",
  "value": "Updated Value",
  "description": "Customer name extracted from the conversation",
  "created_at": 1650000000000,
  "updated_at": 1650000001000
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/audio-to-text",method:"POST",title:"Speech to Text",name:"#audio-to-text"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"This endpoint requires a multipart/form-data request."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsxs(n.p,{children:[`Audio file.
Supported formats: `,e.jsx(n.code,{children:"['mp3', 'mp4', 'mpeg', 'mpga', 'm4a', 'wav', 'webm']"}),`
File size limit: 15MB`]})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"User identifier, defined by the developer's rules, must be unique within the application."})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) Output text"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/audio-to-text",targetCode:`curl -X POST '${i.appDetail.api_base_url}/audio-to-text' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=audio/[mp3|mp4|mpeg|mpga|m4a|wav|webm]'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "text": ""
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"Text to Audio",name:"#text-to-audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Text to speech."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"For text messages generated by Gofy, simply pass the generated message-id directly. The backend will use the message-id to look up the corresponding content and synthesize the voice information directly. If both message_id and text are provided simultaneously, the message_id is given priority."})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"Speech generated content。"})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"The user identifier, defined by the developer, must ensure uniqueness within the app."})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/text-to-audio",targetCode:`curl --location --request POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--form 'text=Hello Gofy;user=abc-123;message_id=5ad4cb98-f0c7-4085-b384-88c403be6290`}),e.jsx(s,{title:"headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"Get Application Basic Information",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get basic information about this application"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) application name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) application description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) application tags"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) application mode"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"author_name"})," (string) application author name"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "advanced-chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"Get Application Parameters Information",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used at the start of entering the page to obtain information such as features, input parameter names, types, and default values."}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"User identifier, defined by the developer's rules, must be unique within the application."})},"user")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) Opening statement"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) List of suggested questions for the opening"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) Suggest questions after enabling the answer.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) Speech to text",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text_to_speech"})," (object) Text to speech",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"voice"})," (string) Voice type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"language"})," (string) Language"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"autoPlay"})," (string) Auto play",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"}),"   Enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"disabled"}),"  Disabled"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) Citation and Attribution",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) Annotation reply",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) User input form configuration",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) Text input control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) Paragraph text input control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) Dropdown control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) Option values"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) File upload configuration",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) Document settings
Currently only supports document types: `,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Document number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) Image settings
Currently only supports image types: `,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Image number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) Audio settings
Currently only supports audio types: `,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Audio number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) Video settings
Currently only supports video types: `,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Video number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) Custom settings",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Custom number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) System parameters",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) Document upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) Image file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) Audio file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) Video file upload size limit (MB)"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "opening_statement": "Hello!",
  "suggested_questions_after_answer": {
      "enabled": true
  },
  "speech_to_text": {
      "enabled": true
  },
  "text_to_speech": {
      "enabled": true,
      "voice": "sambert-zhinan-v1",
      "language": "zh-Hans",
      "autoPlay": "disabled"
  },
  "retriever_resource": {
      "enabled": true
  },
  "annotation_reply": {
      "enabled": true
  },
  "user_input_form": [
      {
          "paragraph": {
              "label": "Query",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/meta",method:"GET",title:"Get Application Meta Information",name:"#meta"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get icons of tools in this application"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_icons"}),"(object[string]) tool icons",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_name"})," (string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (object|string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["(object) icon object",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"background"})," (string) background color in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"content"}),"(string) emoji"]}),`
`]}),`
`]}),`
`,e.jsx(n.li,{children:"(string) url of icon"}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/meta",targetCode:`curl -X GET '${i.appDetail.api_base_url}/meta' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "tool_icons": {
    "dalle2": "https://cloud.gofy.ai/console/api/workspaces/current/tool-provider/builtin/dalle/icon",
    "api_tool": {
      "background": "#252525",
      "content": "\\ud83d\\ude01"
    }
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"Get Application WebApp Settings",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get the WebApp settings of the application."}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) Chat color theme, in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) Whether the chat color theme is inverted"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) Icon type, ",e.jsx(n.code,{children:"emoji"})," - emoji, ",e.jsx(n.code,{children:"image"})," - picture"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) Icon. If it's ",e.jsx(n.code,{children:"emoji"})," type, it's an emoji symbol; if it's ",e.jsx(n.code,{children:"image"})," type, it's an image URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) Background color in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) Icon URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) Description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) Copyright information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) Privacy policy link"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) Custom disclaimer"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) Default language"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) Whether to show workflow details"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) Whether to replace 🤖 in chat with the WebApp icon"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations",method:"GET",title:"Get Annotation List",name:"#annotation_list"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"Page number"})},"page"),e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"Number of items returned, default 20, range 1-100"})},"limit")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/apps/annotations",targetCode:`curl --location --request GET '${i.apiBaseUrl}/apps/annotations?page=1&limit=20' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "data": [
    {
      "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
      "question": "What is your name?",
      "answer": "I am Gofy.",
      "hit_count": 0,
      "created_at": 1735625869
    }
  ],
  "has_more": false,
  "limit": 20,
  "total": 1,
  "page": 1
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations",method:"POST",title:"Create Annotation",name:"#create_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"question",type:"string",children:e.jsx(n.p,{children:"Question"})},"question"),e.jsx(d,{name:"answer",type:"string",children:e.jsx(n.p,{children:"Answer"})},"answer")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/apps/annotations",targetCode:`curl --location --request POST '${i.apiBaseUrl}/apps/annotations' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{"question": "What is your name?","answer": "I am Gofy."}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  {
    "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
    "question": "What is your name?",
    "answer": "I am Gofy.",
    "hit_count": 0,
    "created_at": 1735625869
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations/{annotation_id}",method:"PUT",title:"Update Annotation",name:"#update_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"annotation_id",type:"string",children:e.jsx(n.p,{children:"Annotation ID"})},"annotation_id"),e.jsx(d,{name:"question",type:"string",children:e.jsx(n.p,{children:"Question"})},"question"),e.jsx(d,{name:"answer",type:"string",children:e.jsx(n.p,{children:"Answer"})},"answer")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/apps/annotations/{annotation_id}",targetCode:`curl --location --request POST '${i.apiBaseUrl}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{"question": "What is your name?","answer": "I am Gofy."}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  {
    "id": "69d48372-ad81-4c75-9c46-2ce197b4d402",
    "question": "What is your name?",
    "answer": "I am Gofy.",
    "hit_count": 0,
    "created_at": 1735625869
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotations/{annotation_id}",method:"DELETE",title:"Delete Annotation",name:"#delete_annotation"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"annotation_id",type:"string",children:e.jsx(n.p,{children:"Annotation ID"})},"annotation_id")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/apps/annotations/{annotation_id}",targetCode:`curl --location --request DELETE '${i.apiBaseUrl}/apps/annotations/{annotation_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotation-reply/{action}",method:"POST",title:"Initial Annotation Reply Settings",name:"#initial_annotation_reply_settings"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"action",type:"string",children:e.jsx(n.p,{children:"Action, can only be 'enable' or 'disable'"})},"action"),e.jsx(d,{name:"embedding_model_provider",type:"string",children:e.jsx(n.p,{children:"Specified embedding model provider, must be set up in the system first, corresponding to the provider field(Optional)"})},"embedding_model_provider"),e.jsx(d,{name:"embedding_model",type:"string",children:e.jsx(n.p,{children:"Specified embedding model, corresponding to the model field(Optional)"})},"embedding_model"),e.jsx(d,{name:"score_threshold",type:"number",children:e.jsx(n.p,{children:"The similarity threshold for matching annotated replies. Only annotations with scores above this threshold will be recalled."})},"score_threshold")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.p,{children:"The provider and model name of the embedding model can be obtained through the following interface: v1/workspaces/current/models/model-types/text-embedding. For specific instructions, see: Maintain Knowledge Base via API. The Authorization used is the Dataset API Token."}),e.jsx(s,{title:"Request",tag:"POST",label:"/apps/annotation-reply/{action}",targetCode:`curl --location --request POST '${i.apiBaseUrl}/apps/annotation-reply/{action}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{"score_threshold": 0.9, "embedding_provider_name": "zhipu", "embedding_model_name": "embedding_3"}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "job_id": "b15c8f68-1cf4-4877-bf21-ed7cf2011802",
  "job_status": "waiting"
}
`})})}),e.jsx(n.p,{children:"This interface is executed asynchronously, so it will return a job_id. You can get the final execution result by querying the job status interface."})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/apps/annotation-reply/{action}/status/{job_id}",method:"GET",title:"Query Initial Annotation Reply Settings Task Status",name:"#initial_annotation_reply_settings_task_status"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"action",type:"string",children:e.jsx(n.p,{children:"Action, can only be 'enable' or 'disable', must be the same as the action in the initial annotation reply settings interface"})},"action"),e.jsx(d,{name:"job_id",type:"string",children:e.jsx(n.p,{children:"Job ID,"})},"job_id")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/apps/annotations",targetCode:`curl --location --request GET '${i.apiBaseUrl}/apps/annotation-reply/{action}/status/{job_id}' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "job_id": "b15c8f68-1cf4-4877-bf21-ed7cf2011802",
  "job_status": "waiting",
  "error_msg": ""
}
`})})})]})]})]})}function ge(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(me,{...i})}):me(i)}function be(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"チャットアプリ API"}),`
`,e.jsx(n.p,{children:"チャットアプリケーションはセッションの持続性をサポートしており、以前のチャット履歴を応答のコンテキストとして使用できます。これは、チャットボットやカスタマーサービス AI などに適用できます。"}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"ベース URL"}),e.jsx(s,{title:"コード",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"認証"}),e.jsxs(n.p,{children:["サービス API は ",e.jsx(n.code,{children:"API-Key"}),` 認証を使用します。
`,e.jsx("i",{children:e.jsx(n.strong,{children:"API キーの漏洩を防ぐため、API キーはクライアント側で共有または保存せず、サーバー側で保存することを強くお勧めします。"})})]}),e.jsxs(n.p,{children:["すべての API リクエストにおいて、以下のように ",e.jsx(n.code,{children:"Authorization"}),"HTTP ヘッダーに API キーを含めてください："]}),e.jsx(s,{title:"コード",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages",method:"POST",title:"チャットメッセージを送信",name:"#Send-Chat-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"チャットアプリケーションにリクエストを送信します。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"query",type:"string",children:e.jsx(n.p,{children:"ユーザー入力/質問内容"})},"query"),e.jsx(d,{name:"inputs",type:"object",children:e.jsxs(n.p,{children:[`アプリで定義されたさまざまな変数値の入力を許可します。
`,e.jsx(n.code,{children:"inputs"}),"パラメータには複数のキー/値ペアが含まれ、各キーは特定の変数に対応し、各値はその変数の特定の値です。デフォルトは",e.jsx(n.code,{children:"{}"})]})},"inputs"),e.jsxs(d,{name:"response_mode",type:"string",children:[e.jsx(n.p,{children:"応答の返却モードを指定します。サポートされているモード："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," ストリーミングモード（推奨）、SSE（",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"}),"）を通じてタイプライターのような出力を実装します。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` ブロッキングモード、実行完了後に結果を返します。（プロセスが長い場合、リクエストが中断される可能性があります）
Cloudflareの制限により、100秒後に応答がない場合、リクエストは中断されます。
`,e.jsx("i",{children:"注：エージェントアシスタントモードではブロッキングモードはサポートされていません"})]}),`
`]})]},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーのアイデンティティを定義するために使用されます。
アプリケーション内で開発者によって一意に定義される必要があります。`})},"user"),e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"会話ID、以前のチャット記録に基づいて会話を続けるには、前のメッセージのconversation_idを渡す必要があります。"})},"conversation_id"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"ファイルリスト、モデルが Vision/Video 機能をサポートしている場合に限り、ファイルをテキスト理解および質問応答に組み合わせて入力するのに適しています。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) サポートされるタイプ：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," サポートされるタイプには以下が含まれます：'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," サポートされるタイプには以下が含まれます：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," サポートされるタイプには以下が含まれます：'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," サポートされるタイプには以下が含まれます：'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," サポートされるタイプには以下が含まれます：その他のファイルタイプ"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) 転送方法:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": ファイルのURL。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": ファイルをアップロード。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," ファイルのURL。（転送方法が ",e.jsx(n.code,{children:"remote_url"})," の場合のみ）。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," アップロードされたファイルID。（転送方法が ",e.jsx(n.code,{children:"local_file"})," の場合のみ）。"]}),`
`]})]},"files"),e.jsx(d,{name:"auto_generate_name",type:"bool",children:e.jsxs(n.p,{children:["タイトルを自動生成します。デフォルトは",e.jsx(n.code,{children:"true"}),`です。
`,e.jsx(n.code,{children:"false"}),"に設定すると、会話のリネームAPIを呼び出し、",e.jsx(n.code,{children:"auto_generate"}),"を",e.jsx(n.code,{children:"true"}),"に設定することで非同期タイトル生成を実現できます。"]})},"auto_generate_name"),e.jsx(d,{name:"workflow_id",type:"string",children:e.jsxs(n.p,{children:["（オプション）ワークフローID、特定のバージョンを指定するために使用、提供されない場合はデフォルトの公開バージョンを使用。",e.jsx("br",{}),`
取得方法：バージョン履歴インターフェースで、各バージョンエントリの右側にあるコピーアイコンをクリックすると、完全なワークフローIDをコピーできます。`]})},"workflow_id"),e.jsxs(d,{name:"trace_id",type:"string",children:[e.jsxs(n.p,{children:["（オプション）トレースID。既存の業務システムのトレースコンポーネントと連携し、エンドツーエンドの分散トレーシングを実現するために使用します。指定がない場合、システムが自動的に trace_id を生成します。以下の3つの方法で渡すことができ、優先順位は次のとおりです：",e.jsx("br",{})]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["Header：HTTPヘッダー ",e.jsx("code",{children:"X-Trace-Id"})," で渡す（最優先）。",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["クエリパラメータ：URLクエリパラメータ ",e.jsx("code",{children:"trace_id"})," で渡す。",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["リクエストボディ：リクエストボディの ",e.jsx("code",{children:"trace_id"})," フィールドで渡す（本フィールド）。",e.jsx("br",{})]}),`
`]})]},"trace_id")]}),e.jsx(n.h3,{children:"応答"}),e.jsx(n.p,{children:`response_modeがブロッキングの場合、CompletionResponseオブジェクトを返します。
response_modeがストリーミングの場合、ChunkCompletionResponseストリームを返します。`}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["完全なアプリ結果を返します。",e.jsx(n.code,{children:"Content-Type"}),"は",e.jsx(n.code,{children:"application/json"}),"です。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) イベントタイプ、固定で ",e.jsx(n.code,{children:"message"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ユニークID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) アプリモード、",e.jsx(n.code,{children:"chat"}),"として固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 完全な応答内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) モデル使用情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用と帰属リスト"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) メッセージ作成タイムスタンプ、例：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["アプリによって出力されたストリームチャンクを返します。",e.jsx(n.code,{children:"Content-Type"}),"は",e.jsx(n.code,{children:"text/event-stream"}),`です。
各ストリーミングチャンクは`,e.jsx(n.code,{children:"data:"}),"で始まり、2つの改行文字",e.jsx(n.code,{children:"\\n\\n"}),"で区切られます。以下のように表示されます："]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsxs(n.p,{children:["ストリーミングチャンクの構造は",e.jsx(n.code,{children:"event"}),"に応じて異なります："]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLMはテキストチャンクイベントを返します。つまり、完全なテキストがチャンク形式で出力されます。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLMが返したテキストチャンク内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: agent_message"})," LLMはテキストチャンクイベントを返します。つまり、エージェントアシスタントが有効な場合、完全なテキストがチャンク形式で出力されます（エージェントモードでのみサポート）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLMが返したテキストチャンク内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTSオーディオストリームイベント、つまり音声合成出力。内容はMp3形式のオーディオブロックで、base64文字列としてエンコードされています。再生時には、base64をデコードしてプレーヤーに入力するだけです。（このメッセージは自動再生が有効な場合にのみ利用可能）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 音声合成後のオーディオ、base64テキストコンテンツとしてエンコードされており、再生時にはbase64をデコードしてプレーヤーに入力するだけです"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTSオーディオストリーム終了イベント。このイベントを受信すると、オーディオストリームの終了を示します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 終了イベントにはオーディオがないため、これは空の文字列です"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: agent_thought"})," エージェントの思考、LLMの思考、ツール呼び出しの入力と出力を含みます（エージェントモードでのみサポート）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) エージェント思考ID、各反復には一意のエージェント思考IDがあります"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string)  タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"position"})," (int) 現在のエージェント思考の位置、各メッセージには順番に複数の思考が含まれる場合があります。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"thought"})," (string) LLMが考えていること"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"observation"})," (string) ツール呼び出しからの応答"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool"})," (string) 呼び出されたツールのリスト、;で区切られます"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_input"})," (string) ツールの入力、JSON形式。例：",e.jsx(n.code,{children:'{"dalle3": {"prompt": "a cute cat"}}'}),"。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[string]) message_fileイベントを参照",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"})," (string) ファイルID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_file"})," メッセージファイルイベント、ツールによって新しいファイルが作成されました",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ファイル一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"}),' (string) ファイルタイプ、現在は"image"のみ許可']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) 所属、ここでは'assistant'のみ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) ファイルのリモートURL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"}),"  (string) 会話ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," メッセージ終了イベント、このイベントを受信するとストリーミングが終了したことを意味します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) モデル使用情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用と帰属リスト"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` メッセージ内容置換イベント。
出力内容のモデレーションが有効な場合、内容がフラグされると、このイベントを通じてメッセージ内容が事前設定された返信に置き換えられます。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 置換内容（すべてのLLM返信テキストを直接置換）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
ストリーミングプロセス中に発生した例外はストリームイベントの形式で出力され、エラーイベントを受信するとストリームが終了します。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下のStop Generate APIに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTPステータスコード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) エラーコード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) エラーメッセージ"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 接続を維持するために10秒ごとにpingイベントが発生します。"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404, 会話が存在しません"}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", 異常なパラメータ入力"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", アプリ構成が利用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", 利用可能なモデル資格情報構成がありません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", モデル呼び出しクォータが不足しています"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", 現在のモデルは利用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_not_found"}),", 指定されたワークフローバージョンが見つかりません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"draft_workflow_error"}),", ドラフトワークフローバージョンは使用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_id_format_error"}),", ワークフローID形式エラー、UUID形式が必要です"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"completion_request_error"}),", テキスト生成に失敗しました"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/chat-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "query": "What are the specs of the iPhone 13 Pro Max?",
  "response_mode": "streaming",
  "conversation_id": "",
  "user": "abc-123",
  "files": [
      {
        "type": "image",
      "transfer_method": "remote_url",
      "url": "https://cloud.gofy.ai/logo/logo-site.png"
    }
  ]
}'`}),e.jsx(n.h3,{children:"ブロッキングモード"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "event": "message",
    "task_id": "c3800678-a077-43df-a102-53f23ed20b88",
    "id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "message_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
    "mode": "chat",
    "answer": "iPhone 13 Pro Maxの仕様は次のとおりです:...",
    "metadata": {
        "usage": {
            "prompt_tokens": 1033,
            "prompt_unit_price": "0.001",
            "prompt_price_unit": "0.001",
            "prompt_price": "0.0010330",
            "completion_tokens": 128,
            "completion_unit_price": "0.002",
            "completion_price_unit": "0.001",
            "completion_price": "0.0002560",
            "total_tokens": 1161,
            "total_price": "0.0012890",
            "currency": "USD",
            "latency": 0.7682376249867957
        },
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ]
    },
    "created_at": 1705407629
}
`})})}),e.jsx(n.h3,{children:"ストリーミングモード（基本アシスタント）"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " I", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": "'m", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " glad", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " to", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " meet", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " you", "created_at": 1679586595}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})}),e.jsx(n.h3,{children:"応答例（エージェントアシスタント）"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " I", "created_at": 1679586595}
data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": "'m", "created_at": 1679586595}
data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " glad", "created_at": 1679586595}
data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " to", "created_at": 1679586595}
data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " meet", "created_at": 1679586595}
data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " you", "created_at": 1679586595}
data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"ファイルアップロード",name:"#file-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:`メッセージ送信時に使用するためのファイルをアップロードします（現在は画像のみサポート）。画像とテキストのマルチモーダル理解を可能にします。
png、jpg、jpeg、webp、gif 形式をサポートしています。
アップロードされたファイルは現在のエンドユーザーのみが使用できます。`}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.p,{children:["このインターフェースは",e.jsx(n.code,{children:"multipart/form-data"}),"リクエストを必要とします。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file"}),` (File) 必須
アップロードするファイル。`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子、開発者のルールで定義され、アプリケーション内で一意でなければなりません。サービス API は WebApp によって作成された会話を共有しません。`]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsx(n.p,{children:"アップロードが成功すると、サーバーはファイルの ID と関連情報を返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) ファイル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) ファイルサイズ（バイト）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) ファイル拡張子"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) ファイルの MIME タイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) エンドユーザーID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"no_file_uploaded"}),", ファイルが提供されなければなりません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"too_many_files"}),", 現在は 1 つのファイルのみ受け付けます"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_preview"}),", ファイルはプレビューをサポートしていません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_estimate"}),", ファイルは推定をサポートしていません"]}),`
`,e.jsxs(n.li,{children:["413, ",e.jsx(n.code,{children:"file_too_large"}),", ファイルが大きすぎます"]}),`
`,e.jsxs(n.li,{children:["415, ",e.jsx(n.code,{children:"unsupported_file_type"}),", サポートされていない拡張子、現在はドキュメントファイルのみ受け付けます"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_connection_failed"}),", S3 サービスに接続できません"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_permission_denied"}),", S3 にファイルをアップロードする権限がありません"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_file_too_large"}),", ファイルが S3 のサイズ制限を超えています"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"エンドユーザーを取得",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"エンドユーザー ID からエンドユーザー情報を取得します。"}),e.jsxs(n.p,{children:["他の API がエンドユーザー ID（例：ファイルアップロードの ",e.jsx(n.code,{children:"created_by"}),"）を返す場合に利用できます。"]}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) 必須
エンドユーザー ID。`]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"EndUser オブジェクトを返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) テナント ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) アプリ ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) エンドユーザー種別"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) 外部ユーザー ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) 匿名ユーザーかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) セッション ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 日時"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 日時"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"end_user_not_found"}),", エンドユーザーが見つかりません"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"レスポンス例"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"ファイルプレビュー",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アップロードされたファイルをプレビューまたはダウンロードします。このエンドポイントを使用すると、以前にファイルアップロード API でアップロードされたファイルにアクセスできます。"}),e.jsx("i",{children:"ファイルは、リクエストしているアプリケーションのメッセージ範囲内にある場合のみアクセス可能です。"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) 必須
プレビューするファイルの一意識別子。ファイルアップロード API レスポンスから取得します。`]}),`
`]}),e.jsx(n.h3,{children:"クエリパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) オプション
ファイルを添付ファイルとして強制ダウンロードするかどうか。デフォルトは `,e.jsx(n.code,{children:"false"}),"（ブラウザでプレビュー）。"]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"ブラウザ表示またはダウンロード用の適切なヘッダー付きでファイル内容を返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," ファイル MIME タイプに基づいて設定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," ファイルサイズ（バイト、利用可能な場合）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"})," ",e.jsx(n.code,{children:"as_attachment=true"}),' の場合は "attachment" に設定']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," パフォーマンス向上のためのキャッシュヘッダー"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' 音声/動画ファイルの場合は "bytes" に設定']}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", パラメータ入力異常"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", ファイルアクセス拒否またはファイルが現在のアプリケーションに属していません"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", ファイルが見つからないか削除されています"]}),`
`,e.jsx(n.li,{children:"500, サーバー内部エラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"添付ファイルとしてダウンロード"}),e.jsx(s,{title:"Download Request",tag:"GET",label:"/files/:file_id/preview?as_attachment=true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"レスポンスヘッダー例"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"ダウンロードレスポンスヘッダー"}),e.jsx(s,{title:"Download Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages/:task_id/stop",method:"POST",title:"生成停止",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ストリーミングモードでのみサポートされています。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、ストリーミングチャンクの返り値から取得できます"]}),`
`]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子、エンドユーザーのアイデンティティを定義するために使用され、メッセージ送信インターフェースで渡されたユーザーと一致している必要があります。サービス API は WebApp によって作成された会話を共有しません。`]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"POST",label:"/chat-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{"user": "abc-123"}'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"メッセージフィードバック",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"エンドユーザーはフィードバックメッセージを提供でき、アプリケーション開発者が期待される出力を最適化するのに役立ちます。"}),e.jsx(n.h3,{children:"パス"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"メッセージID"})},"message_id")}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsxs(n.p,{children:["アップボートは",e.jsx(n.code,{children:"like"}),"、ダウンボートは",e.jsx(n.code,{children:"dislike"}),"、アップボートの取り消しは",e.jsx(n.code,{children:"null"})]})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者のルールで定義され、アプリケーション内で一意でなければなりません。"})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"メッセージのフィードバックです。"})},"content")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"アプリのメッセージの「いいね」とフィードバックを取得",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アプリのエンドユーザーからのフィードバックや「いいね」を取得します。"}),e.jsx(n.h3,{children:"クエリ"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（任意）ページ番号。デフォルト値：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（任意）1ページあたりの件数。デフォルト値：20"})},"limit")}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (リスト) このアプリの「いいね」とフィードバックの一覧を返します。"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`    {
    "data": [
        {
            "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
            "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
            "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
            "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
            "rating": "like",
            "content": "message feedback information-3",
            "from_source": "user",
            "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
            "from_account_id": null,
            "created_at": "2025-04-24T09:24:38",
            "updated_at": "2025-04-24T09:24:38"
        }
    ]
    }
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/{message_id}/suggested",method:"GET",title:"次の推奨質問",name:"#suggested"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"現在のメッセージに対する次の質問の提案を取得します"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"メッセージID"})},"message_id")}),e.jsx(n.h3,{children:"クエリ"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーのアイデンティティを定義するために使用され、統計のために使用されます。
アプリケーション内で開発者によって一意に定義される必要があります。`})},"user")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/messages/{message_id}/suggested",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/messages/{message_id}/suggested?user=abc-123& \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success",
  "data": [
        "a",
        "b",
        "c"
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages",method:"GET",title:"会話履歴メッセージを取得",name:"#messages"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:["スクロールロード形式で過去のチャット記録を返し、最初のページは最新の",e.jsx(n.code,{children:"{limit}"}),"メッセージを返します。つまり、逆順です。"]}),e.jsx(n.h3,{children:"クエリ"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"会話ID"})},"conversation_id"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーのアイデンティティを定義するために使用され、統計のために使用されます。
アプリケーション内で開発者によって一意に定義される必要があります。`})},"user"),e.jsx(d,{name:"first_id",type:"string",children:e.jsx(n.p,{children:"現在のページの最初のチャット記録のID、デフォルトはnullです。"})},"first_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"1回のリクエストで返すチャット履歴メッセージの数、デフォルトは20です。"})},"limit")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) メッセージリスト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) メッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ユーザー入力パラメータ。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"})," (string) ユーザー入力/質問内容。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[object]) メッセージファイル",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) ファイルタイプ、画像の場合はimage"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) ファイルプレビューURL、ファイルアクセスにはファイルプレビューAPI（",e.jsx(n.code,{children:"/files/{file_id}/preview"}),"）を使用してください"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) 所属、ユーザーまたはアシスタント"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"agent_thoughts"})," (array[object]) エージェントの思考（基本アシスタントの場合は空）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) エージェント思考ID、各反復には一意のエージェント思考IDがあります"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"position"})," (int) 現在のエージェント思考の位置、各メッセージには順番に複数の思考が含まれる場合があります。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"thought"})," (string) LLMが考えていること"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"observation"})," (string) ツール呼び出しからの応答"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool"})," (string) 呼び出されたツールのリスト、;で区切られます"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_input"})," (string) ツールの入力、JSON形式。例：",e.jsx(n.code,{children:'{"dalle3": {"prompt": "a cute cat"}}'}),"。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[string]) message_fileイベントを参照",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"})," (string) ファイルID"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 応答メッセージ内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"feedback"})," (object) フィードバック情報",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"rating"})," (string) アップボートは",e.jsx(n.code,{children:"like"})," / ダウンボートは",e.jsx(n.code,{children:"dislike"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用と帰属リスト"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) 次のページがあるかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返されたアイテムの数、入力がシステム制限を超える場合、システム制限の数を返します"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/messages",targetCode:`curl -X GET '${i.appDetail.api_base_url}/messages?user=abc-123&conversation_id='\\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"応答例（基本アシスタント）"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
        "id": "a076a87f-31e5-48dc-b452-0061adbbc922",
        "conversation_id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
        "inputs": {
            "name": "gofy"
        },
        "query": "iphone 13 pro",
        "answer": "iPhone 13 Proは2021年9月24日に発売され、6.1インチのディスプレイと1170 x 2532の解像度を備えています。Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)プロセッサ、6 GBのRAMを搭載し、128 GB、256 GB、512 GB、1 TBのストレージオプションを提供します。カメラは12 MP、バッテリー容量は3095 mAhで、iOS 15を搭載しています。",
        "message_files": [],
        "feedback": null,
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ],
        "agent_thoughts": [],
        "created_at": 1705569239,
    }
  ]
}
`})})}),e.jsx(n.h3,{children:"応答例（エージェントアシスタント）"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "limit": 20,
    "has_more": false,
    "data": [
        {
            "id": "d35e006c-7c4d-458f-9142-be4930abdf94",
            "conversation_id": "957c068b-f258-4f89-ba10-6e8a0361c457",
            "inputs": {},
            "query": "draw a cat",
            "answer": "猫の画像を生成しました。メッセージを確認して画像を表示してください。",
            "message_files": [
                {
                    "id": "976990d2-5294-47e6-8f14-7356ba9d2d76",
                    "type": "image",
                    "url": "http://127.0.0.1:5001/files/tools/976990d2-5294-47e6-8f14-7356ba9d2d76.png?timestamp=1705988524&nonce=55df3f9f7311a9acd91bf074cd524092&sign=z43nMSO1L2HBvoqADLkRxr7Biz0fkjeDstnJiCK1zh8=",
                    "belongs_to": "assistant"
                }
            ],
            "feedback": null,
            "retriever_resources": [],
            "created_at": 1705988187,
            "agent_thoughts": [
                {
                    "id": "592c84cf-07ee-441c-9dcc-ffc66c033469",
                    "chain_id": null,
                    "message_id": "d35e006c-7c4d-458f-9142-be4930abdf94",
                    "position": 1,
                    "thought": "",
                    "tool": "dalle2",
                    "tool_input": "{\\"dalle2\\": {\\"prompt\\": \\"cat\\"}}",
                    "created_at": 1705988186,
                    "observation": "画像はすでに作成され、ユーザーに送信されました。今すぐユーザーに確認するように伝えてください。",
                    "files": [
                        "976990d2-5294-47e6-8f14-7356ba9d2d76"
                    ]
                },
                {
                    "id": "73ead60d-2370-4780-b5ed-532d2762b0e5",
                    "chain_id": null,
                    "message_id": "d35e006c-7c4d-458f-9142-be4930abdf94",
                    "position": 2,
                    "thought": "猫の画像を生成しました。メッセージを確認して画像を表示してください。",
                    "tool": "",
                    "tool_input": "",
                    "created_at": 1705988199,
                    "observation": "",
                    "files": []
                }
            ]
        }
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations",method:"GET",title:"会話を取得",name:"#conversations"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"現在のユーザーの会話リストを取得し、デフォルトで最新の 20 件を返します。"}),e.jsx(n.h3,{children:"クエリ"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`ユーザー識別子、エンドユーザーのアイデンティティを定義するために使用され、統計のために使用されます。
アプリケーション内で開発者によって一意に定義される必要があります。`})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"(Optional)現在のページの最後のレコードのID、デフォルトはnullです。"})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"(Optional)1回のリクエストで返すレコードの数、デフォルトは最新の20件です。最大100、最小1。"})},"limit"),e.jsxs(d,{name:"sort_by",type:"string",children:[e.jsx(n.p,{children:"(Optional)ソートフィールド、デフォルト：-updated_at（更新時間で降順にソート）"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"利用可能な値：created_at, -created_at, updated_at, -updated_at"}),`
`,e.jsx(n.li,{children:'フィールドの前の記号は順序または逆順を表し、"-"は逆順を表します。'}),`
`]})]},"sort_by")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 会話のリスト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 会話名、デフォルトでは、ユーザーが会話で最初に尋ねた質問のスニペットです。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ユーザー入力パラメータ。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) 紹介"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) 更新タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返されたエントリの数、入力がシステム制限を超える場合、システム制限の数を返します"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/conversations",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations?user=abc-123&last_id=&limit=20' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
      "id": "10799fb8-64f7-4296-bbf7-b42bfbe0ae54",
      "name": "新しいチャット",
      "inputs": {
          "book": "book",
          "myName": "Lucy"
      },
      "status": "normal",
      "created_at": 1679667915,
      "updated_at": 1679667915
    },
    {
      "id": "hSIhXBhNe8X1d8Et"
      // ...
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id",method:"DELETE",title:"会話を削除",name:"#delete"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"会話を削除します。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者によって定義され、アプリケーション内で一意である必要があります。"})},"user")}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"DELETE",label:"/conversations/:conversation_id",targetCode:`curl -X DELETE '${i.appDetail.api_base_url}/conversations/:conversation_id' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "user": "abc-123"
}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/name",method:"POST",title:"会話の名前を変更",name:"#rename"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"リクエストボディ"}),e.jsx(n.p,{children:"セッションの名前を変更します。セッション名は、複数のセッションをサポートするクライアントでの表示に使用されます。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会話ID"]}),`
`]}),e.jsxs(t,{children:[e.jsx(d,{name:"name",type:"string",children:e.jsxs(n.p,{children:["(Optional)会話の名前。このパラメータは、",e.jsx(n.code,{children:"auto_generate"}),"が",e.jsx(n.code,{children:"true"}),"に設定されている場合、省略できます。"]})},"name"),e.jsx(d,{name:"auto_generate",type:"bool",children:e.jsxs(n.p,{children:["(Optional)タイトルを自動生成します。デフォルトは",e.jsx(n.code,{children:"false"}),"です。"]})},"auto_generate"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者によって定義され、アプリケーション内で一意である必要があります。"})},"user")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 会話ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 会話名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ユーザー入力パラメータ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 会話状態"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) 紹介"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) 更新タイムスタンプ、例：1705395332"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/conversations/:conversation_id/name",targetCode:`curl -X POST '${i.appDetail.api_base_url}/conversations/:conversation_id/name' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "name": "",
  "auto_generate": true,
  "user": "abc-123"
}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
    "name": "Chat vs AI",
    "inputs": {},
    "introduction": "",
    "created_at": 1705569238,
    "updated_at": 1705569238
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables",method:"GET",title:"会話変数の取得",name:"#conversation-variables"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"特定の会話から変数を取得します。このエンドポイントは、会話中に取得された構造化データを抽出するのに役立ちます。"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsx(t,{children:e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"変数を取得する会話のID。"})},"conversation_id")}),e.jsx(n.h3,{children:"クエリパラメータ"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子。開発者によって定義されたルールに従い、アプリケーション内で一意である必要があります。"})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"(Optional)現在のページの最後のレコードのID、デフォルトはnullです。"})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"(Optional)1回のリクエストで返すレコードの数、デフォルトは最新の20件です。最大100、最小1。"})},"limit")]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) ページごとのアイテム数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) さらにアイテムがあるかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 変数のリスト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 変数 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 変数名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) 変数タイプ（文字列、数値、真偽値など）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (string) 変数値"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 変数の説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) 最終更新タイムスタンプ"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", 会話が見つかりません"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/conversations/:conversation_id/variables",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Request with variable name filter",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X GET '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123&variable_name=customer_name' \\
--header 'Authorization: Bearer {api_key}'
`})})}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 100,
  "has_more": false,
  "data": [
    {
      "id": "variable-uuid-1",
      "name": "customer_name",
      "value_type": "string",
      "value": "John Doe",
      "description": "会話から抽出された顧客名",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    },
    {
      "id": "variable-uuid-2",
      "name": "order_details",
      "value_type": "json",
      "value": "{\\"product\\":\\"Widget\\",\\"quantity\\":5,\\"price\\":19.99}",
      "description": "顧客の注文詳細",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables/:variable_id",method:"PUT",title:"会話変数の更新",name:"#update-conversation-variable"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"特定の会話変数の値を更新します。このエンドポイントは、名前、型、説明を保持しながら、会話中にキャプチャされた変数の値を変更することを可能にします。"}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"更新する変数を含む会話のID。"})},"conversation_id"),e.jsx(d,{name:"variable_id",type:"string",children:e.jsx(n.p,{children:"更新する変数のID。"})},"variable_id")]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"value",type:"any",children:e.jsx(n.p,{children:"変数の新しい値。変数の期待される型（文字列、数値、オブジェクトなど）と一致する必要があります。"})},"value"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子。開発者によって定義されたルールに従い、アプリケーション内で一意である必要があります。"})},"user")]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"以下を含む更新された変数オブジェクトを返します："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 変数名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) 変数型（文字列、数値、オブジェクトなど）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (any) 更新された変数値"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 変数の説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) 最終更新タイムスタンプ"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"Type mismatch: variable expects {expected_type}, but got {actual_type} type"}),", 値の型が変数の期待される型と一致しません"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", 会話が見つかりません"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_variable_not_exists"}),", 変数が見つかりません"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"PUT",label:"/conversations/:conversation_id/variables/:variable_id",targetCode:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "value": "Updated Value",
  "user": "abc-123"
}'`}),e.jsxs(s,{title:"異なる値型での更新",children:[e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": "新しい文字列値",
    "user": "abc-123"
}'
`})}),e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": 42,
    "user": "abc-123"
}'
`})}),e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": {"product": "Widget", "quantity": 10, "price": 29.99},
    "user": "abc-123"
}'
`})})]}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "variable-uuid-1",
  "name": "customer_name",
  "value_type": "string",
  "value": "Updated Value",
  "description": "会話から抽出された顧客名",
  "created_at": 1650000000000,
  "updated_at": 1650000001000
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/audio-to-text",method:"POST",title:"音声からテキストへ",name:"#audio-to-text"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"このエンドポイントは multipart/form-data リクエストを必要とします。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsxs(n.p,{children:[`オーディオファイル。
サポートされている形式：`,e.jsx(n.code,{children:"['mp3', 'mp4', 'mpeg', 'mpga', 'm4a', 'wav', 'webm']"}),`
ファイルサイズ制限：15MB`]})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者のルールで定義され、アプリケーション内で一意でなければなりません。"})},"user")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) 出力テキスト"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/audio-to-text",targetCode:`curl -X POST '${i.appDetail.api_base_url}/audio-to-text' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=audio/[mp3|mp4|mpeg|mpga|m4a|wav|webm]'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "text": ""
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"テキストから音声へ",name:"#text-to-audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"テキストを音声に変換します。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"Gofyによって生成されたテキストメッセージの場合、生成されたメッセージIDを直接渡します。バックエンドはメッセージIDを使用して対応するコンテンツを検索し、音声情報を直接合成します。message_idとtextが同時に提供される場合、message_idが優先されます。"})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"音声生成コンテンツ。"})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"ユーザー識別子、開発者によって定義され、アプリ内で一意である必要があります。"})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/text-to-audio",targetCode:`curl --location --request POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--form 'text=Hello Gofy;user=abc-123;message_id=5ad4cb98-f0c7-4085-b384-88c403be6290`}),e.jsx(s,{title:"ヘッダー",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"アプリケーションの基本情報を取得",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"このアプリケーションの基本情報を取得するために使用されます"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) アプリケーションの名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) アプリケーションの説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) アプリケーションのタグ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) アプリケーションのモード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"author_name"})," (string) 作者の名前"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"アプリケーションのパラメータ情報を取得",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ページに入る際に、機能、入力パラメータ名、タイプ、デフォルト値などの情報を取得するために使用されます。"}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) 開始文"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) 開始時の推奨質問のリスト"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) 答えを有効にした後の質問を提案します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) 音声からテキストへ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text_to_speech"})," (object) テキストから音声へ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"voice"})," (string) 音声タイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"language"})," (string) 言語"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"autoPlay"})," (string) 自動再生",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"}),"  有効"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"disabled"})," 無効"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) 引用と帰属",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) 注釈返信",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) ユーザー入力フォームの構成",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) テキスト入力コントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) 段落テキスト入力コントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) ドロップダウンコントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) オプション値"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) ファイルアップロード設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) ドキュメント設定
現在サポートされているドキュメントタイプ：`,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) ドキュメント数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) 画像設定
現在サポートされている画像タイプ：`,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 画像数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) オーディオ設定
現在サポートされているオーディオタイプ：`,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) オーディオ数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) ビデオ設定
現在サポートされているビデオタイプ：`,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) ビデオ数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) カスタム設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) カスタム数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) システムパラメータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) ドキュメントアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) 画像ファイルアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) オーディオファイルアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) ビデオファイルアップロードサイズ制限（MB）"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "opening_statement": "こんにちは！",
  "suggested_questions_after_answer": {
      "enabled": true
  },
  "speech_to_text": {
      "enabled": true
  },
  "text_to_speech": {
      "enabled": true,
      "voice": "sambert-zhinan-v1",
      "language": "zh-Hans",
      "autoPlay": "disabled"
  },
  "retriever_resource": {
      "enabled": true
  },
  "annotation_reply": {
      "enabled": true
  },
  "user_input_form": [
      {
          "paragraph": {
              "label": "クエリ",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/meta",method:"GET",title:"アプリケーションのメタ情報を取得",name:"#meta"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"このアプリケーションのツールのアイコンを取得するために使用されます"}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_icons"}),"(object[string]) ツールアイコン",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_name"})," (string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (object|string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["(object) アイコンオブジェクト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"background"})," (string) 背景色（16 進数形式）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"content"}),"(string) 絵文字"]}),`
`]}),`
`]}),`
`,e.jsx(n.li,{children:"(string) アイコンの URL"}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/meta",targetCode:`curl -X GET '${i.appDetail.api_base_url}/meta' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "tool_icons": {
    "dalle2": "https://cloud.gofy.ai/console/api/workspaces/current/tool-provider/builtin/dalle/icon",
    "api_tool": {
      "background": "#252525",
      "content": "\\ud83d\\ude01"
    }
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"アプリのWebApp設定を取得",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アプリの WebApp 設定を取得するために使用します。"}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp 名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) チャットの色テーマ、16 進数形式"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) チャットの色テーマを反転するかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) アイコンタイプ、",e.jsx(n.code,{children:"emoji"}),"-絵文字、",e.jsx(n.code,{children:"image"}),"-画像"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) アイコン。",e.jsx(n.code,{children:"emoji"}),"タイプの場合は絵文字、",e.jsx(n.code,{children:"image"}),"タイプの場合は画像 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) 16 進数形式の背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) アイコンの URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) 著作権情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) プライバシーポリシーのリンク"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) カスタム免責事項"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) デフォルト言語"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) ワークフローの詳細を表示するかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) WebApp のアイコンをチャット内の🤖に置き換えるかどうか"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{})]})}function fe(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(be,{...i})}):be(i)}function qe(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"对话型应用 API"}),`
`,e.jsx(n.p,{children:"对话应用支持会话持久化，可将之前的聊天记录作为上下文进行回答，可适用于聊天/客服 AI 等。"}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"基础 URL"}),e.jsx(s,{title:"Code",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"鉴权"}),e.jsxs(n.p,{children:["Service API 使用 ",e.jsx(n.code,{children:"API-Key"}),` 进行鉴权。
`,e.jsx("i",{children:e.jsxs(n.strong,{children:["强烈建议开发者把 ",e.jsx(n.code,{children:"API-Key"})," 放在后端存储，而非分享或者放在客户端存储，以免 ",e.jsx(n.code,{children:"API-Key"})," 泄露，导致财产损失。"]})}),`
所有 API 请求都应在 `,e.jsx(n.strong,{children:e.jsx(n.code,{children:"Authorization"})})," HTTP Header 中包含您的 ",e.jsx(n.code,{children:"API-Key"}),"，如下所示："]}),e.jsx(s,{title:"Code",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages",method:"POST",title:"发送对话消息",name:"#Create-Chat-Message"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"创建会话消息。"}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"query",type:"string",children:e.jsx(n.p,{children:"用户输入/提问内容。"})},"query"),e.jsx(d,{name:"inputs",type:"object",children:e.jsxs(n.p,{children:[`允许传入 App 定义的各变量值。
inputs 参数包含了多组键值对（Key/Value pairs），每组的键对应一个特定变量，每组的值则是该变量的具体值。
默认 `,e.jsx(n.code,{children:"{}"})]})},"inputs"),e.jsx(d,{name:"response_mode",type:"string",children:e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," 流式模式（推荐）。基于 SSE（",e.jsx(n.strong,{children:e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"})}),"）实现类似打字机输出方式的流式返回。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` 阻塞模式，等待执行完毕后返回结果。（请求若流程较长可能会被中断）。
`,e.jsx("i",{children:"由于 Cloudflare 限制，请求会在 100 秒超时无返回后中断。"}),`
注：Agent模式下不允许blocking。`]}),`
`]})},"response_mode"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:`用户标识，用于定义终端用户的身份，方便检索、统计。
由开发者定义规则，需保证用户标识在应用内唯一。服务 API 不会共享 WebApp 创建的对话。`})},"user"),e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"（选填）会话 ID，需要基于之前的聊天记录继续对话，必须传之前消息的 conversation_id。"})},"conversation_id"),e.jsxs(d,{name:"files",type:"array[object]",children:[e.jsx(n.p,{children:"文件列表，适用于传入文件结合文本理解并回答问题，仅当模型支持 Vision/Video 能力时可用。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 支持类型：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," 具体类型包含：'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," 具体类型包含：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," 具体类型包含：'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," 具体类型包含：'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," 具体类型包含：其他文件类型"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string)  传递方式:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": 文件地址。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": 上传文件。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," 文件地址。（仅当传递方式为 ",e.jsx(n.code,{children:"remote_url"})," 时）。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," 上传文件 ID。（仅当传递方式为 ",e.jsx(n.code,{children:"local_file "}),"时）。"]}),`
`]})]},"files"),e.jsx(d,{name:"auto_generate_name",type:"bool",children:e.jsxs(n.p,{children:["（选填）自动生成标题，默认 ",e.jsx(n.code,{children:"true"}),"。 若设置为 ",e.jsx(n.code,{children:"false"}),"，则可通过调用会话重命名接口并设置 ",e.jsx(n.code,{children:"auto_generate"})," 为 ",e.jsx(n.code,{children:"true"})," 实现异步生成标题。"]})},"auto_generate_name"),e.jsx(d,{name:"workflow_id",type:"string",children:e.jsxs(n.p,{children:["（选填）工作流ID，用于指定特定版本，如果不提供则使用默认的已发布版本。",e.jsx("br",{}),`
获取方式：在版本历史界面，点击每个版本条目右侧的复制图标即可复制完整的工作流 ID。`]})},"workflow_id"),e.jsxs(d,{name:"trace_id",type:"string",children:[e.jsxs(n.p,{children:["（选填）链路追踪ID。适用于与业务系统已有的trace组件打通，实现端到端分布式追踪等场景。如果未指定，系统会自动生成",e.jsx("code",{children:"trace_id"}),"。支持以下三种方式传递，具体优先级依次为：",e.jsx("br",{})]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["Header：通过 HTTP Header ",e.jsx("code",{children:"X-Trace-Id"})," 传递，优先级最高。",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["Query 参数：通过 URL 查询参数 ",e.jsx("code",{children:"trace_id"})," 传递。",e.jsx("br",{})]}),`
`,e.jsxs(n.li,{children:["Request Body：通过请求体字段 ",e.jsx("code",{children:"trace_id"})," 传递（即本字段）。",e.jsx("br",{})]}),`
`]})]},"trace_id")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(t,{children:[e.jsxs(n.p,{children:["当 ",e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"blocking"}),` 时，返回 ChatCompletionResponse object。
当 `,e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"streaming"}),"时，返回 ChunkChatCompletionResponse object 流式序列。"]}),e.jsx(n.h3,{children:"ChatCompletionResponse"}),e.jsxs(n.p,{children:["返回完整的 App 结果，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"application/json"}),"。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 事件类型，固定为 ",e.jsx(n.code,{children:"message"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 唯一ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) App 模式，固定为 chat"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 完整回复内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) 模型用量信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用和归属分段列表"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 消息创建时间戳，如：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"ChunkChatCompletionResponse"}),e.jsxs(n.p,{children:["返回 App 输出的流式块，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"text/event-stream"}),`。
每个流式块均为 data: 开头，块之间以 \\n\\n 即两个换行符分隔，如下所示：`]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "message", "task_id": "900bbd43-dc0b-4383-a372-aa6e6c414227", "id": "663c5084-a254-4040-8ad3-51f2a3c1a77c", "answer": "Hi", "created_at": 1705398420}\\n\\n
`})})}),e.jsx(n.p,{children:"流式块中根据 event 不同，结构也不同："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message"})," LLM 返回文本块事件，即：完整的文本以分块的方式输出。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLM 返回文本块内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: agent_message"})," Agent模式下返回文本块事件，即：在Agent模式下，文章的文本以分块的方式输出（仅Agent模式下使用）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) LLM 返回文本块内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: agent_thought"})," Agent模式下有关Agent思考步骤的相关内容，涉及到工具调用（仅Agent模式下使用）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) agent_thought ID，每一轮Agent迭代都会有一个唯一的id"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务ID，用于请求跟踪下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"position"})," (int) agent_thought在消息中的位置，如第一轮迭代position为1"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"thought"})," (string) agent的思考内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"observation"})," (string) 工具调用的返回结果"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool"})," (string) 使用的工具列表，以 ; 分割多个工具"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_input"})," (string) 工具的输入，JSON格式的字符串(object)。如：",e.jsx(n.code,{children:'{"dalle3": {"prompt": "a cute cat"}}'})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[string])  当前 ",e.jsx(n.code,{children:"agent_thought"})," 关联的文件ID",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"})," (string) 文件ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_file"})," 文件事件，表示有新文件需要展示",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 文件唯一ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 文件类型，目前仅为image"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) 文件归属，user或assistant，该接口返回仅为 ",e.jsx(n.code,{children:"assistant"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) 文件访问地址"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"}),"  (string) 会话ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_end"})," 消息结束事件，收到此事件则代表流式返回结束。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"metadata"})," (object) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"usage"})," (Usage) 模型用量信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用和归属分段列表"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS 音频流事件，即：语音合成输出。内容是Mp3格式的音频块，使用 base64 编码后的字符串，播放的时候直接解码即可。(开启自动播放才有此消息)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 语音合成之后的音频块使用 Base64 编码之后的文本内容，播放的时候直接 base64 解码送入播放器即可"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS 音频流结束事件，收到这个事件表示音频流返回结束。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 结束事件是没有音频的，所以这里是空字符串"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: message_replace"}),` 消息内容替换事件。
开启内容审查和审查输出内容时，若命中了审查条件，则会通过此事件替换消息内容为预设回复。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string) 替换内容（直接替换 LLM 所有回复文本）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: error"}),`
流式输出过程中出现的异常会以 stream event 形式输出，收到异常事件后即结束。`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (int) HTTP 状态码"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"code"})," (string) 错误码"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message"})," (string) 错误消息"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 每 10s 一次的 ping 事件，保持连接存活。"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"404，对话不存在"}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"invalid_param"}),"，传入参数异常"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"app_unavailable"}),"，App 配置不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_not_initialize"}),"，无可用模型凭据配置"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_quota_exceeded"}),"，模型调用额度不足"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"model_currently_not_support"}),"，当前模型不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_not_found"}),"，指定的工作流版本未找到"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"draft_workflow_error"}),"，无法使用草稿工作流版本"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_id_format_error"}),"，工作流ID格式错误，需要UUID格式"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"completion_request_error"}),"，文本生成失败"]}),`
`,e.jsx(n.li,{children:"500，服务内部异常"}),`
`]})]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/chat-messages",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "query": "What are the specs of the iPhone 13 Pro Max?",
  "response_mode": "streaming",
  "conversation_id": "",
  "user": "abc-123",
  "files": [
      {
        "type": "image",
      "transfer_method": "remote_url",
      "url": "https://cloud.gofy.ai/logo/logo-site.png"
    }
  ]
}'`}),e.jsx(n.h3,{children:"阻塞模式"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "event": "message",
    "task_id": "c3800678-a077-43df-a102-53f23ed20b88",
    "id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "message_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2",
    "mode": "chat",
    "answer": "iPhone 13 Pro Max specs are listed here:...",
    "metadata": {
        "usage": {
            "prompt_tokens": 1033,
            "prompt_unit_price": "0.001",
            "prompt_price_unit": "0.001",
            "prompt_price": "0.0010330",
            "completion_tokens": 128,
            "completion_unit_price": "0.002",
            "completion_price_unit": "0.001",
            "completion_price": "0.0002560",
            "total_tokens": 1161,
            "total_price": "0.0012890",
            "currency": "USD",
            "latency": 0.7682376249867957
        },
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ]
    },
    "created_at": 1705407629
}
`})})}),e.jsx(n.h3,{children:"流式模式（基础助手）"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " I", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": "'m", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " glad", "created_at": 1679586595}
  data: {"event": "message", "message_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " to", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " meet", "created_at": 1679586595}
  data: {"event": "message", "message_id" : "5ad4cb98-f0c7-4085-b384-88c403be6290", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "answer": " you", "created_at": 1679586595}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})}),e.jsx(n.h3,{children:"流式模式（智能助手）"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "agent_thought", "id": "8dcf3648-fbad-407a-85dd-73a6f43aeb9f", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "position": 1, "thought": "", "observation": "", "tool": "", "tool_input": "", "created_at": 1705639511, "message_files": [], "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "agent_thought", "id": "8dcf3648-fbad-407a-85dd-73a6f43aeb9f", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "position": 1, "thought": "", "observation": "", "tool": "dalle3", "tool_input": "{\\"dalle3\\": {\\"prompt\\": \\"cute Japanese anime girl with white hair, blue eyes, bunny girl suit\\"}}", "created_at": 1705639511, "message_files": [], "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "message_file", "id": "d75b7a5c-ce5e-442e-ab1b-d6a5e5b557b0", "type": "image", "belongs_to": "assistant", "url": "http://127.0.0.1:5001/files/tools/d75b7a5c-ce5e-442e-ab1b-d6a5e5b557b0.png?timestamp=1705639526&nonce=70423256c60da73a9c96d1385ff78487&sign=7B5fKV9890YJuqchQvrABvW4AIupDvDvxGdu1EOJT94=", "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "agent_thought", "id": "8dcf3648-fbad-407a-85dd-73a6f43aeb9f", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "position": 1, "thought": "", "observation": "image has been created and sent to user already, you should tell user to check it now.", "tool": "dalle3", "tool_input": "{\\"dalle3\\": {\\"prompt\\": \\"cute Japanese anime girl with white hair, blue eyes, bunny girl suit\\"}}", "created_at": 1705639511, "message_files": ["d75b7a5c-ce5e-442e-ab1b-d6a5e5b557b0"], "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "agent_thought", "id": "67a99dc1-4f82-42d3-b354-18d4594840c8", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "position": 2, "thought": "", "observation": "", "tool": "", "tool_input": "", "created_at": 1705639511, "message_files": [], "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "agent_message", "id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "answer": "I have created an image of a cute Japanese", "created_at": 1705639511, "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "agent_message", "id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "answer": " anime girl with white hair and blue", "created_at": 1705639511, "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "agent_message", "id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "answer": " eyes wearing a bunny girl" ,"created_at": 1705639511, "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "agent_message", "id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "answer": " suit .", "created_at": 1705639511, "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "agent_thought", "id": "67a99dc1-4f82-42d3-b354-18d4594840c8", "task_id": "9cf1ddd7-f94b-459b-b942-b77b26c59e9b", "message_id": "1fb10045-55fd-4040-99e6-d048d07cbad3", "position": 2, "thought": "I have created an image of a cute Japanese anime girl with white hair and blue eyes wearing a bunny girl suit.", "observation": "", "tool": "", "tool_input": "", "created_at": 1705639511, "message_files": [], "conversation_id": "c216c595-2d89-438c-b33c-aae5ddddd142"}
  data: {"event": "message_end", "id": "5e52ce04-874b-4d27-9045-b3bc80def685", "conversation_id": "45701982-8118-4bc5-8e9b-64562b4555f2", "metadata": {"usage": {"prompt_tokens": 1033, "prompt_unit_price": "0.001", "prompt_price_unit": "0.001", "prompt_price": "0.0010330", "completion_tokens": 135, "completion_unit_price": "0.002", "completion_price_unit": "0.001", "completion_price": "0.0002700", "total_tokens": 1168, "total_price": "0.0013030", "currency": "USD", "latency": 1.381760165997548}, "retriever_resources": [{"position": 1, "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb", "dataset_name": "iPhone", "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00", "document_name": "iPhone List", "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a", "score": 0.98457545, "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""}]}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"上传文件",name:"#files-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:[`上传文件（目前仅支持图片）并在发送消息时使用，可实现图文多模态理解。
支持 png, jpg, jpeg, webp, gif 格式。
`,e.jsx("i",{children:"上传的文件仅供当前终端用户使用。"})]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["该接口需使用  ",e.jsx(n.code,{children:"multipart/form-data"})," 进行请求。"]}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsx(n.p,{children:"要上传的文件。"})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，用于定义终端用户的身份，必须和发送消息接口传入 user 保持一致。服务 API 不会共享 WebApp 创建的对话。"})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"成功上传后，服务器会返回文件的 ID 和相关信息。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 文件名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) 文件大小（byte）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) 文件后缀"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) 文件 mime-type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) 上传人 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 上传时间"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"no_file_uploaded"}),"，必须提供文件"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"too_many_files"}),"，目前只接受一个文件"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"unsupported_preview"}),"，该文件不支持预览"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"unsupported_estimate"}),"，该文件不支持估算"]}),`
`,e.jsxs(n.li,{children:["413，",e.jsx(n.code,{children:"file_too_large"}),"，文件太大"]}),`
`,e.jsxs(n.li,{children:["415，",e.jsx(n.code,{children:"unsupported_file_type"}),"，不支持的扩展名，当前只接受文档类文件"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_connection_failed"}),"，无法连接到 S3 服务"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_permission_denied"}),"，无权限上传文件到 S3"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_file_too_large"}),"，文件超出 S3 大小限制"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": 123,
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"获取终端用户",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"通过终端用户 ID 获取终端用户信息。"}),e.jsxs(n.p,{children:["当其他 API 返回终端用户 ID（例如：上传文件接口返回的 ",e.jsx(n.code,{children:"created_by"}),"）时，可使用该接口查询对应的终端用户信息。"]}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) 必需
终端用户 ID。`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"返回 EndUser 对象。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) 工作空间（Tenant）ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) 应用 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 终端用户类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) 外部用户 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) 是否匿名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 时间"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404，",e.jsx(n.code,{children:"end_user_not_found"}),"，终端用户不存在"]}),`
`,e.jsx(n.li,{children:"500，内部服务器错误"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/:file_id/preview",method:"GET",title:"文件预览",name:"#file-preview"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"预览或下载已上传的文件。此端点允许您访问先前通过文件上传 API 上传的文件。"}),e.jsx("i",{children:"文件只能在属于请求应用程序的消息范围内访问。"}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"}),` (string) 必需
要预览的文件的唯一标识符，从文件上传 API 响应中获得。`]}),`
`]}),e.jsx(n.h3,{children:"查询参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"as_attachment"}),` (boolean) 可选
是否强制将文件作为附件下载。默认为 `,e.jsx(n.code,{children:"false"}),"（在浏览器中预览）。"]}),`
`]}),e.jsx(n.h3,{children:"响应"}),e.jsx(n.p,{children:"返回带有适当浏览器显示或下载标头的文件内容。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Type"})," 根据文件 MIME 类型设置"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Length"})," 文件大小（以字节为单位，如果可用）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Content-Disposition"})," 如果 ",e.jsx(n.code,{children:"as_attachment=true"}),' 则设置为 "attachment"']}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Cache-Control"})," 用于性能的缓存标头"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"Accept-Ranges"}),' 对于音频/视频文件设置为 "bytes"']}),`
`]}),e.jsx(n.h3,{children:"错误"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", 参数输入异常"]}),`
`,e.jsxs(n.li,{children:["403, ",e.jsx(n.code,{children:"file_access_denied"}),", 文件访问被拒绝或文件不属于当前应用程序"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"file_not_found"}),", 文件未找到或已被删除"]}),`
`,e.jsx(n.li,{children:"500, 服务内部错误"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"请求示例"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'`,children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X GET '\${props.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview' \\
--header 'Authorization: Bearer {api_key}'
`})})}),e.jsx(n.h3,{children:"作为附件下载"}),e.jsx(s,{title:"Request",tag:"GET",label:"/files/:file_id/preview?as_attachment=true",targetCode:`curl -X GET '${i.appDetail.api_base_url}/files/72fa9618-8f89-4a37-9b33-7e1178a24a67/preview?as_attachment=true' \\
--header 'Authorization: Bearer {api_key}' \\
--output downloaded_file.png`}),e.jsx(n.h3,{children:"响应标头示例"}),e.jsx(s,{title:"Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Cache-Control: public, max-age=3600
`})})}),e.jsx(n.h3,{children:"文件下载响应标头"}),e.jsx(s,{title:"Download Response Headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-http",children:`Content-Type: image/png
Content-Length: 1024
Content-Disposition: attachment; filename*=UTF-8''example.png
Cache-Control: public, max-age=3600
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/chat-messages/:task_id/stop",method:"POST",title:"停止响应",name:"#Stop"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"仅支持流式模式。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，可在流式返回 Chunk 中获取"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
用户标识，用于定义终端用户的身份，必须和发送消息接口传入 user 保持一致。API 无法访问 WebApp 创建的会话。`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"})," (string) 固定返回 success"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/chat-messages/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/chat-messages/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{ "user": "abc-123"}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/:message_id/feedbacks",method:"POST",title:"消息反馈（点赞）",name:"#feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"消息终端用户反馈、点赞，方便应用开发者优化输出预期。"}),e.jsx(n.h3,{children:"Path Params"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"消息 ID"})},"message_id")}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"rating",type:"string",children:e.jsx(n.p,{children:"点赞 like, 点踩 dislike,  撤销点赞 null"})},"rating"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。服务 API 不会共享 WebApp 创建的对话。"})},"user"),e.jsx(d,{name:"content",type:"string",children:e.jsx(n.p,{children:"消息反馈的具体信息。"})},"content")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"})," (string) 固定返回 success"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/messages/:message_id/feedbacks",targetCode:`curl -X POST '${i.appDetail.api_base_url}/messages/:message_id/feedbacks \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "rating": "like",
  "user": "abc-123",
  "content": "message feedback information"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/app/feedbacks",method:"GET",title:"获取APP的消息点赞和反馈",name:"#app-feedbacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"获取应用的终端用户反馈、点赞。"}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"page",type:"string",children:e.jsx(n.p,{children:"（选填）分页，默认值：1"})},"page")}),e.jsx(t,{children:e.jsx(d,{name:"limit",type:"string",children:e.jsx(n.p,{children:"（选填）每页数量，默认值：20"})},"limit")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (List) 返回该APP的点赞、反馈列表。"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/app/feedbacks",targetCode:`curl -X GET '${i.appDetail.api_base_url}/app/feedbacks?page=1&limit=20'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`    {
    "data": [
        {
            "id": "8c0fbed8-e2f9-49ff-9f0e-15a35bdd0e25",
            "app_id": "f252d396-fe48-450e-94ec-e184218e7346",
            "conversation_id": "2397604b-9deb-430e-b285-4726e51fd62d",
            "message_id": "709c0b0f-0a96-4a4e-91a4-ec0889937b11",
            "rating": "like",
            "content": "message feedback information-3",
            "from_source": "user",
            "from_end_user_id": "74286412-9a1a-42c1-929c-01edb1d381d5",
            "from_account_id": null,
            "created_at": "2025-04-24T09:24:38",
            "updated_at": "2025-04-24T09:24:38"
        }
    ]
    }
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages/{message_id}/suggested",method:"GET",title:"获取下一轮建议问题列表",name:"#suggested"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"获取下一轮建议问题列表。"}),e.jsx(n.h3,{children:"Path Params"}),e.jsx(t,{children:e.jsx(d,{name:"message_id",type:"string",children:e.jsx(n.p,{children:"Message ID"})},"message_id")}),e.jsx(n.h3,{children:"Query"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/messages/{message_id}/suggested",targetCode:`curl --location --request GET '${i.appDetail.api_base_url}/messages/{message_id}/suggested?user=abc-123 \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--header 'Content-Type: application/json'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success",
  "data": [
        "a",
        "b",
        "c"
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/messages",method:"GET",title:"获取会话历史消息",name:"#messages"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:["滚动加载形式返回历史聊天记录，第一页返回最新  ",e.jsx(n.code,{children:"limit"})," 条，即：倒序返回。"]}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"会话 ID"})},"conversation_id"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user"),e.jsx(d,{name:"first_id",type:"string",children:e.jsx(n.p,{children:"当前页第一条聊天记录的 ID，默认 null"})},"first_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"一次请求返回多少条聊天记录，默认 20 条。"})},"limit")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object])  消息列表",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"}),"  (string) 消息 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string)  会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 用户输入参数。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"query"}),"  (string) 用户输入 / 提问内容。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[object]) 消息文件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 文件类型，image 图片"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) 文件预览地址，使用文件预览 API (",e.jsx(n.code,{children:"/files/{file_id}/preview"}),") 访问文件"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"belongs_to"})," (string) 文件归属方，user 或 assistant"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"agent_thoughts"})," (array[object]) Agent思考内容（仅Agent模式下不为空）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) agent_thought ID，每一轮Agent迭代都会有一个唯一的id"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"position"})," (int) agent_thought在消息中的位置，如第一轮迭代position为1"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"thought"})," (string) agent的思考内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"observation"})," (string) 工具调用的返回结果"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool"})," (string) 使用的工具列表，以 ; 分割多个工具"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_input"})," (string) 工具的输入，JSON格式的字符串(object)。如：",e.jsx(n.code,{children:'{"dalle3": {"prompt": "a cute cat"}}'})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_files"})," (array[string])  当前agent_thought 关联的文件ID",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_id"})," (string) 文件ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话ID"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"answer"})," (string)  回答消息内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"}),"  (timestamp) 创建时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"feedback"})," (object) 反馈信息",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"rating"})," (string) 点赞 like / 点踩 dislike"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resources"})," (array[RetrieverResource]) 引用和归属分段列表"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) 是否存在下一页"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返回条数，若传入超过系统限制，返回系统限制数量"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/messages",targetCode:`curl -X GET '${i.appDetail.api_base_url}/messages?user=abc-123&conversation_id=' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example(基础助手)"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
"limit": 20,
"has_more": false,
"data": [
    {
        "id": "a076a87f-31e5-48dc-b452-0061adbbc922",
        "conversation_id": "cd78daf6-f9e4-4463-9ff2-54257230a0ce",
        "inputs": {
            "name": "gofy"
        },
        "query": "iphone 13 pro",
        "answer": "The iPhone 13 Pro, released on September 24, 2021, features a 6.1-inch display with a resolution of 1170 x 2532. It is equipped with a Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard) processor, 6 GB of RAM, and offers storage options of 128 GB, 256 GB, 512 GB, and 1 TB. The camera is 12 MP, the battery capacity is 3095 mAh, and it runs on iOS 15.",
        "message_files": [],
        "feedback": null,
        "retriever_resources": [
            {
                "position": 1,
                "dataset_id": "101b4c97-fc2e-463c-90b1-5261a4cdcafb",
                "dataset_name": "iPhone",
                "document_id": "8dd1ad74-0b5f-4175-b735-7d98bbbb4e00",
                "document_name": "iPhone List",
                "segment_id": "ed599c7f-2766-4294-9d1d-e5235a61270a",
                "score": 0.98457545,
                "content": "\\"Model\\",\\"Release Date\\",\\"Display Size\\",\\"Resolution\\",\\"Processor\\",\\"RAM\\",\\"Storage\\",\\"Camera\\",\\"Battery\\",\\"Operating System\\"\\n\\"iPhone 13 Pro Max\\",\\"September 24, 2021\\",\\"6.7 inch\\",\\"1284 x 2778\\",\\"Hexa-core (2x3.23 GHz Avalanche + 4x1.82 GHz Blizzard)\\",\\"6 GB\\",\\"128, 256, 512 GB, 1TB\\",\\"12 MP\\",\\"4352 mAh\\",\\"iOS 15\\""
            }
        ],
        "agent_thoughts": [],
        "created_at": 1705569239
    }
  ]
}
`})})}),e.jsx(n.h3,{children:"Response Example(智能助手)"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
"limit": 20,
"has_more": false,
"data": [
    {
        "id": "d35e006c-7c4d-458f-9142-be4930abdf94",
        "conversation_id": "957c068b-f258-4f89-ba10-6e8a0361c457",
        "inputs": {},
        "query": "draw a cat",
        "answer": "I have generated an image of a cat for you. Please check your messages to view the image.",
        "message_files": [
            {
                "id": "976990d2-5294-47e6-8f14-7356ba9d2d76",
                "type": "image",
                "url": "http://127.0.0.1:5001/files/tools/976990d2-5294-47e6-8f14-7356ba9d2d76.png?timestamp=1705988524&nonce=55df3f9f7311a9acd91bf074cd524092&sign=z43nMSO1L2HBvoqADLkRxr7Biz0fkjeDstnJiCK1zh8=",
                "belongs_to": "assistant"
            }
        ],
        "feedback": null,
        "retriever_resources": [],
        "created_at": 1705988187,
        "agent_thoughts": [
            {
                "id": "592c84cf-07ee-441c-9dcc-ffc66c033469",
                "chain_id": null,
                "message_id": "d35e006c-7c4d-458f-9142-be4930abdf94",
                "position": 1,
                "thought": "",
                "tool": "dalle2",
                "tool_input": "{\\"dalle2\\": {\\"prompt\\": \\"cat\\"}}",
                "created_at": 1705988186,
                "observation": "image has been created and sent to user already, you should tell user to check it now.",
                "files": [
                    "976990d2-5294-47e6-8f14-7356ba9d2d76"
                ]
            },
            {
                "id": "73ead60d-2370-4780-b5ed-532d2762b0e5",
                "chain_id": null,
                "message_id": "d35e006c-7c4d-458f-9142-be4930abdf94",
                "position": 2,
                "thought": "I have generated an image of a cat for you. Please check your messages to view the image.",
                "tool": "",
                "tool_input": "",
                "created_at": 1705988199,
                "observation": "",
                "files": []
            }
        ]
    }
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations",method:"GET",title:"获取会话列表",name:"#conversations"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"获取当前用户的会话列表，默认返回最近的 20 条。"}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"（选填）当前页最后面一条记录的 ID，默认 null"})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"（选填）一次请求返回多少条记录，默认 20 条，最大 100 条，最小 1 条。"})},"limit"),e.jsxs(d,{name:"sort_by",type:"string",children:[e.jsx(n.p,{children:"（选填）排序字段，默认 -updated_at(按更新时间倒序排列)"}),e.jsxs(n.ul,{children:[`
`,e.jsx(n.li,{children:"可选值：created_at, -created_at, updated_at, -updated_at"}),`
`,e.jsx(n.li,{children:"字段前面的符号代表顺序或倒序，-代表倒序"}),`
`]})]},"sort_by")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 会话列表",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"}),"  (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"}),"  (string) 会话名称，默认为会话中用户最开始问题的截取。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 用户输入参数。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 会话状态"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) 开场白"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 创建时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) 更新时间"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返回条数，若传入超过系统限制，返回系统限制数量"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/conversations",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations?user=abc-123&last_id=&limit=20'\\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 20,
  "has_more": false,
  "data": [
    {
      "id": "10799fb8-64f7-4296-bbf7-b42bfbe0ae54",
      "name": "New chat",
      "inputs": {
          "book": "book",
          "myName": "Lucy"
      },
      "status": "normal",
      "created_at": 1679667915,
      "updated_at": 1679667915
    },
    {
      "id": "hSIhXBhNe8X1d8Et"
      // ...
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id",method:"DELETE",title:"删除会话",name:"#delete"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"删除会话。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsx(t,{children:e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"})," (string) 固定返回 success"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"DELETE",label:"/conversations/:conversation_id",targetCode:`curl -X DELETE '${i.appDetail.api_base_url}/conversations/:conversation_id' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-text",children:`204 No Content
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/name",method:"POST",title:"会话重命名",name:"#rename"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"对会话进行重命名，会话名称用于显示在支持多会话的客户端上。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"conversation_id"})," (string) 会话 ID"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"name",type:"string",children:e.jsxs(n.p,{children:["（选填）名称，若 ",e.jsx(n.code,{children:"auto_generate"})," 为 ",e.jsx(n.code,{children:"true"})," 时，该参数可不传。"]})},"name"),e.jsx(d,{name:"auto_generate",type:"bool",children:e.jsx(n.p,{children:"（选填）自动生成标题，默认 false。"})},"auto_generate"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"}),"  (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"}),"  (string) 会话名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 用户输入参数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 会话状态"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"introduction"})," (string) 开场白"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 创建时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (timestamp) 更新时间"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/conversations/:conversation_id/name",targetCode:`curl -X POST '${i.appDetail.api_base_url}/conversations/:conversation_id/name' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "name": "",
  "auto_generate": true,
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "34d511d5-56de-4f16-a997-57b379508443",
  "name": "hello",
  "inputs": {},
  "status": "normal",
  "introduction": "",
  "created_at": 1732731141,
  "updated_at": 1732734510
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables",method:"GET",title:"获取对话变量",name:"#conversation-variables"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"从特定对话中检索变量。此端点对于提取对话过程中捕获的结构化数据非常有用。"}),e.jsx(n.h3,{children:"路径参数"}),e.jsx(t,{children:e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"要从中检索变量的对话ID。"})},"conversation_id")}),e.jsx(n.h3,{children:"查询参数"}),e.jsxs(t,{children:[e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识符，由开发人员定义的规则，在应用程序内必须唯一。"})},"user"),e.jsx(d,{name:"last_id",type:"string",children:e.jsx(n.p,{children:"（选填）当前页最后面一条记录的 ID，默认 null"})},"last_id"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"（选填）一次请求返回多少条记录，默认 20 条，最大 100 条，最小 1 条。"})},"limit")]}),e.jsx(n.h3,{children:"响应"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 每页项目数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) 是否有更多项目"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 变量列表",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 变量 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 变量名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) 变量类型（字符串、数字、布尔等）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (string) 变量值"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 变量描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) 最后更新时间戳"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"错误"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", 对话不存在"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/conversations/:conversation_id/variables",targetCode:`curl -X GET '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Request with variable name filter",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X GET '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables?user=abc-123&variable_name=customer_name' \\
--header 'Authorization: Bearer {api_key}'
`})})}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "limit": 100,
  "has_more": false,
  "data": [
    {
      "id": "variable-uuid-1",
      "name": "customer_name",
      "value_type": "string",
      "value": "John Doe",
      "description": "客户名称（从对话中提取）",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    },
    {
      "id": "variable-uuid-2",
      "name": "order_details",
      "value_type": "json",
      "value": "{\\"product\\":\\"Widget\\",\\"quantity\\":5,\\"price\\":19.99}",
      "description": "客户的订单详情",
      "created_at": 1650000000000,
      "updated_at": 1650000000000
    }
  ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/conversations/:conversation_id/variables/:variable_id",method:"PUT",title:"更新对话变量",name:"#update-conversation-variable"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"更新特定对话变量的值。此端点允许您修改在对话过程中捕获的变量值，同时保留其名称、类型和描述。"}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(t,{children:[e.jsx(d,{name:"conversation_id",type:"string",children:e.jsx(n.p,{children:"包含要更新变量的对话ID。"})},"conversation_id"),e.jsx(d,{name:"variable_id",type:"string",children:e.jsx(n.p,{children:"要更新的变量ID。"})},"variable_id")]}),e.jsx(n.h3,{children:"请求体"}),e.jsxs(t,{children:[e.jsx(d,{name:"value",type:"any",children:e.jsx(n.p,{children:"变量的新值。必须匹配变量的预期类型（字符串、数字、对象等）。"})},"value"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识符，由开发人员定义的规则，在应用程序内必须唯一。"})},"user")]}),e.jsx(n.h3,{children:"响应"}),e.jsx(n.p,{children:"返回包含以下内容的更新变量对象："}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 变量ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 变量名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value_type"})," (string) 变量类型（字符串、数字、对象等）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"value"})," (any) 更新后的变量值"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 变量描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (int) 最后更新时间戳"]}),`
`]}),e.jsx(n.h3,{children:"错误"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"Type mismatch: variable expects {expected_type}, but got {actual_type} type"}),", 值类型与变量的预期类型不匹配"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_not_exists"}),", 对话不存在"]}),`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"conversation_variable_not_exists"}),", 变量不存在"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"PUT",label:"/conversations/:conversation_id/variables/:variable_id",targetCode:`curl -X PUT '${i.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "value": "Updated Value",
  "user": "abc-123"
}'`}),e.jsxs(s,{title:"使用不同值类型更新",children:[e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": "新的字符串值",
    "user": "abc-123"
}'
`})}),e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": 42,
    "user": "abc-123"
}'
`})}),e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-bash",children:`curl -X PUT '\${props.appDetail.api_base_url}/conversations/{conversation_id}/variables/{variable_id}' \\
--header 'Content-Type: application/json' \\
--header 'Authorization: Bearer {api_key}' \\
--data-raw '{
    "value": {"product": "Widget", "quantity": 10, "price": 29.99},
    "user": "abc-123"
}'
`})})]}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "variable-uuid-1",
  "name": "customer_name",
  "value_type": "string",
  "value": "Updated Value",
  "description": "客户名称（从对话中提取）",
  "created_at": 1650000000000,
  "updated_at": 1650000001000
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/audio-to-text",method:"POST",title:"语音转文字",name:"#audio-to-text"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["该接口需使用 ",e.jsx(n.code,{children:"multipart/form-data"})," 进行请求。"]}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsxs(n.p,{children:[`语音文件。
支持格式：`,e.jsx(n.code,{children:"['mp3', 'mp4', 'mpeg', 'mpga', 'm4a', 'wav', 'webm']"}),`
文件大小限制：15MB`]})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) 输出文字"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/audio-to-text",targetCode:`curl -X POST '${i.appDetail.api_base_url}/audio-to-text' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=audio/[mp3|mp4|mpeg|mpga|m4a|wav|webm]`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "text": "hello"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/text-to-audio",method:"POST",title:"文字转语音",name:"#text-to-audio"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"文字转语音。"}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(t,{children:[e.jsx(d,{name:"message_id",type:"str",children:e.jsx(n.p,{children:"Gofy 生成的文本消息，那么直接传递生成的message-id 即可，后台会通过 message_id 查找相应的内容直接合成语音信息。如果同时传 message_id 和 text，优先使用 message_id。"})},"message_id"),e.jsx(d,{name:"text",type:"str",children:e.jsx(n.p,{children:"语音生成内容。如果没有传 message-id的话，则会使用这个字段的内容"})},"text"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，由开发者定义规则，需保证用户标识在应用内唯一。"})},"user")]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/text-to-audio",targetCode:`curl --location --request POST '${i.appDetail.api_base_url}/text-to-audio' \\
--header 'Authorization: Bearer ENTER-YOUR-SECRET-KEY' \\
--form 'text=你好Gofy;user=abc-123;message_id=5ad4cb98-f0c7-4085-b384-88c403be6290`}),e.jsx(s,{title:"headers",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "Content-Type": "audio/wav"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"获取应用基本信息",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取应用的基本信息"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 应用名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 应用描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) 应用标签"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) 应用模式"]}),`
`,e.jsx(n.li,{children:"'author_name' (string) 作者名称"}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "chat",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"获取应用参数",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于进入页面一开始，获取功能开关、输入参数名称、类型及默认值等使用。"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"opening_statement"})," (string) 开场白"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions"})," (array[string]) 开场推荐问题列表"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"suggested_questions_after_answer"})," (object) 启用回答后给出推荐问题。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"speech_to_text"})," (object) 语音转文本",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text_to_speech"})," (object) 文本转语音",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"voice"})," (string) 语音类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"language"})," (string) 语言"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"autoPlay"})," (string) 自动播放",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"}),"  开启"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"disabled"})," 关闭"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"retriever_resource"})," (object) 引用和归属",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"annotation_reply"})," (object) 标记回复",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否开启"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) 用户输入表单配置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) 文本输入控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) 段落文本输入控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) 下拉控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) 选项值"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) 文件上传配置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) 文档设置
当前仅支持文档类型：`,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 文档数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) 图片设置
当前仅支持图片类型：`,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 图片数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) 音频设置
当前仅支持音频类型：`,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 音频数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) 视频设置
当前仅支持视频类型：`,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 视频数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) 自定义设置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 自定义数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) 系统参数",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) 文档上传大小限制 (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) 图片文件上传大小限制（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) 音频文件上传大小限制 (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) 视频文件上传大小限制 (MB)"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'\\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "introduction": "nice to meet you",
  "user_input_form": [
    {
      "text-input": {
        "label": "a",
        "variable": "a",
        "required": true,
        "max_length": 48,
        "default": ""
      }
    },
    {
      // ...
    }
  ],
  "file_upload": {
    "image": {
      "enabled": true,
      "number_limits": 3,
      "transfer_methods": [
        "remote_url",
        "local_file"
      ]
    }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/meta",method:"GET",title:"获取应用Meta信息",name:"#meta"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取工具 icon"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tool_icons"}),"(object[string]) 工具图标",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"工具名称"})," (string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (object|string)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["(object) 图标",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"background"})," (string) hex 格式的背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"content"}),"(string) emoji"]}),`
`]}),`
`]}),`
`,e.jsx(n.li,{children:"(string) 图标 URL"}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"POST",label:"/meta",targetCode:`curl -X GET '${i.appDetail.api_base_url}/meta' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "tool_icons": {
      "dalle2": "https://cloud.gofy.ai/console/api/workspaces/current/tool-provider/builtin/dalle/icon",
      "api_tool": {
          "background": "#252525",
          "content": "\\ud83d\\ude01"
      }
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"获取应用 WebApp 设置",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取应用的 WebApp 设置"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp 名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme"})," (string) 聊天颜色主题，hex 格式"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"chat_color_theme_inverted"})," (bool) 聊天颜色主题是否反转"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) 图标类型，",e.jsx(n.code,{children:"emoji"}),"-表情，",e.jsx(n.code,{children:"image"}),"-图片"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) 图标，如果是 ",e.jsx(n.code,{children:"emoji"})," 类型，则是 emoji 表情符号，如果是 ",e.jsx(n.code,{children:"image"})," 类型，则是图片 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) hex 格式的背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) 图标 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) 版权信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) 隐私政策链接"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) 自定义免责声明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) 默认语言"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) 是否显示工作流详情"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"use_icon_as_answer_icon"})," (bool) 是否使用 WebApp 图标替换聊天中的 🤖"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "chat_color_theme": "#ff4a4a",
  "chat_color_theme_inverted": false,
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
  "use_icon_as_answer_icon": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{})]})}function ye(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(qe,{...i})}):qe(i)}function ve(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"Workflow App API"}),`
`,e.jsx(n.p,{children:"Workflow applications offers non-session support and is ideal for translation, article writing, summarization AI, and more."}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"Base URL"}),e.jsx(s,{title:"Code",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"Authentication"}),e.jsxs(n.p,{children:["The Service API uses ",e.jsx(n.code,{children:"API-Key"}),` authentication.
`,e.jsx("i",{children:e.jsx(n.strong,{children:"Strongly recommend storing your API Key on the server-side, not shared or stored on the client-side, to avoid possible API-Key leakage that can lead to serious consequences."})})]}),e.jsxs(n.p,{children:["For all API requests, include your API Key in the ",e.jsx(n.code,{children:"Authorization"})," HTTP Header, as shown below:"]}),e.jsx(s,{title:"Code",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/run",method:"POST",title:"Execute Workflow",name:"#Execute-Workflow"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Execute workflow, cannot be executed without a published workflow."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"inputs"}),` (object) Required
Allows the entry of various variable values defined by the App.
The `,e.jsx(n.code,{children:"inputs"}),` parameter contains multiple key/value pairs, with each key corresponding to a specific variable and each value being the specific value for that variable.
The workflow application requires at least one key/value pair to be inputted. The variable can be of File Array type.
File Array type variable is suitable for inputting files combined with text understanding and answering questions, available only when the model supports file parsing and understanding capability.
If the variable is of File Array type, the corresponding value should be a list whose elements contain following attributions:`]}),`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) Supported type:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," Supported types include: 'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," Supported types include: 'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," Supported types include: 'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," Supported types include: 'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," Supported types include: other file types"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) Transfer method:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": File URL."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": Upload file."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," File URL. (Only when transfer method is ",e.jsx(n.code,{children:"remote_url"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," Upload file ID. (Only when transfer method is ",e.jsx(n.code,{children:"local_file"}),")."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"response_mode"}),` (string) Required
The mode of response return, supporting:`]}),`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," Streaming mode (recommended), implements a typewriter-like output through SSE (",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` Blocking mode, returns result after execution is complete. (Requests may be interrupted if the process is long)
`,e.jsx("i",{children:"Due to Cloudflare restrictions, the request will be interrupted without a return after 100 seconds."})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`]}),`
`,e.jsx("br",{}),`
`,e.jsx("i",{children:"The user identifier should be consistent with the user passed in the message sending interface. The Service API does not share conversations created by the WebApp."}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"files"})," (array[object]) Optional"]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"trace_id"}),` (string) Optional
Trace ID. Used for integration with existing business trace components to achieve end-to-end distributed tracing. If not provided, the system will automatically generate a trace_id. Supports the following three ways to pass, in order of priority:`]}),`
`,e.jsxs(n.ol,{children:[`
`,e.jsxs(n.li,{children:["Header: via HTTP Header ",e.jsx(n.code,{children:"X-Trace-Id"}),", highest priority."]}),`
`,e.jsxs(n.li,{children:["Query parameter: via URL query parameter ",e.jsx(n.code,{children:"trace_id"}),"."]}),`
`,e.jsxs(n.li,{children:["Request Body: via request body field ",e.jsx(n.code,{children:"trace_id"})," (i.e., this field)."]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.p,{children:["When ",e.jsx(n.code,{children:"response_mode"})," is ",e.jsx(n.code,{children:"blocking"}),`, return a CompletionResponse object.
When `,e.jsx(n.code,{children:"response_mode"})," is ",e.jsx(n.code,{children:"streaming"}),", return a ChunkCompletionResponse stream."]}),e.jsx(n.h3,{children:"CompletionResponse"}),e.jsxs(n.p,{children:["Returns the App result, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"application/json"}),"."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail of result",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) default 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) start time"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) end time"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"ChunkCompletionResponse"}),e.jsxs(n.p,{children:["Returns the stream chunks outputted by the App, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"text/event-stream"}),`.
Each streaming chunk starts with `,e.jsx(n.code,{children:"data:"}),", separated by two newline characters ",e.jsx(n.code,{children:"\\n\\n"}),", as shown below:"]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "text_chunk", "workflow_run_id": "b85e5fc5-751b-454d-b14e-dc5f240b0a31", "task_id": "bd029338-b068-4d34-a331-fc85478922c2", "data": {"text": "\\u4e3a\\u4e86", "from_variable_selector": ["1745912968134", "text"]}}\\n\\n
`})})}),e.jsxs(n.p,{children:["The structure of the streaming chunks varies depending on the ",e.jsx(n.code,{children:"event"}),":"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," workflow starts execution",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"workflow_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," node execution started",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"node_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ID of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) type of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) name of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) Execution sequence number, used to display Tracing Node sequence"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional Prefix node ID, used for canvas display execution path"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) Contents of all preceding node variables used in the node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) timestamp of start, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: text_chunk"})," Text fragment",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"text_chunk"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) Text content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"from_variable_selector"})," (array) Text source path, helping developers understand which node and variable generated the text"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," node execution ends, success or failure in different states in the same event",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"node_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ID of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) type of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) name of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) Execution sequence number, used to display Tracing Node sequence"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional Prefix node ID, used for canvas display execution path"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) Contents of all preceding node variables used in the node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) Optional node process data"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) meta data",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) optional tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) optional Total cost"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) optional e.g. ",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) timestamp of start, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," workflow execution ends, success or failure in different states in the same event",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"workflow_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) default 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) start time"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) end time"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS audio stream event, that is, speech synthesis output. The content is an audio block in Mp3 format, encoded as a base64 string. When playing, simply decode the base64 and feed it into the player. (This message is available only when auto-play is enabled)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the stop response interface below"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The audio after speech synthesis, encoded in base64 text content, when playing, simply decode the base64 and feed it into the player"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g.: 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS audio stream end event, receiving this event indicates the end of the audio stream.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the stop response interface below"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The end event has no audio, so this is an empty string"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g.: 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," Ping event every 10 seconds to keep the connection alive."]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", abnormal parameter input"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", App configuration unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", no available model credential configuration"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", model invocation quota insufficient"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", current model unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_request_error"}),", workflow execution failed"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/workflows/run",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/run' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Example: file array as an input variable",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "inputs": {
    "{variable_name}":
    [
      {
      "transfer_method": "local_file",
      "upload_file_id": "{upload_file_id}",
      "type": "{document_type}"
      }
    ]
  }
}
`})})}),e.jsx(n.h3,{children:"Blocking Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "workflow_run_id": "djflajgkldjgd",
    "task_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "data": {
        "id": "fdlsjfjejkghjda",
        "workflow_id": "fldjaslkfjlsda",
        "status": "succeeded",
        "outputs": {
          "text": "Nice to meet you."
        },
        "error": null,
        "elapsed_time": 0.875,
        "total_tokens": 3562,
        "total_steps": 8,
        "created_at": 1705407629,
        "finished_at": 1727807631
    }
}
`})})}),e.jsx(n.h3,{children:"Streaming Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})}),e.jsx(s,{title:"File upload sample code",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`import requests
import json

def upload_file(file_path, user):
    upload_url = "https://api.gofy.ai/v1/files/upload"
    headers = {
        "Authorization": "Bearer app-xxxxxxxx",
    }

    try:
        print("Upload file...")
        with open(file_path, 'rb') as file:
            files = {
                'file': (file_path, file, 'text/plain')  # Make sure the file is uploaded with the appropriate MIME type
            }
            data = {
                "user": user,
                "type": "TXT"  # Set the file type to TXT
            }

            response = requests.post(upload_url, headers=headers, files=files, data=data)
            if response.status_code == 201:  # 201 means creation is successful
                print("File uploaded successfully")
                return response.json().get("id")  # Get the uploaded file ID
            else:
                print(f"File upload failed, status code: {response.status_code}")
                return None
    except Exception as e:
        print(f"Error occurred: {str(e)}")
        return None

def run_workflow(file_id, user, response_mode="blocking"):
    workflow_url = "https://api.gofy.ai/v1/workflows/run"
    headers = {
        "Authorization": "Bearer app-xxxxxxxxx",
        "Content-Type": "application/json"
    }

    data = {
        "inputs": {
            "orig_mail": [{
                "transfer_method": "local_file",
                "upload_file_id": file_id,
                "type": "document"
            }]
        },
        "response_mode": response_mode,
        "user": user
    }

    try:
        print("Run Workflow...")
        response = requests.post(workflow_url, headers=headers, json=data)
        if response.status_code == 200:
            print("Workflow execution successful")
            return response.json()
        else:
            print(f"Workflow execution failed, status code: {response.status_code}")
            return {"status": "error", "message": f"Failed to execute workflow, status code: {response.status_code}"}
    except Exception as e:
        print(f"Error occurred: {str(e)}")
        return {"status": "error", "message": str(e)}

# Usage Examples
file_path = "{your_file_path}"
user = "gofyuser"

# Upload files
file_id = upload_file(file_path, user)
if file_id:
    # The file was uploaded successfully, and the workflow continues to run
    result = run_workflow(file_id, user)
    print(result)
else:
    print("File upload failed and workflow cannot be executed")
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/:workflow_id/run",method:"POST",title:"Execute Specific Workflow",name:"#Execute-Specific-Workflow"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Execute a specific version of workflow by specifying the workflow ID in the path parameter."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) Required Workflow ID to specify a specific version of workflow"]}),`
`]}),e.jsx(n.p,{children:"How to obtain: In the version history interface, click the copy icon on the right side of each version entry to copy the complete workflow ID. Each version entry contains a copyable ID field."}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"inputs"}),` (object) Required
Allows the entry of various variable values defined by the App.
The `,e.jsx(n.code,{children:"inputs"}),` parameter contains multiple key/value pairs, with each key corresponding to a specific variable and each value being the specific value for that variable.
The workflow application requires at least one key/value pair to be inputted. The variable can be of File Array type.
File Array type variable is suitable for inputting files combined with text understanding and answering questions, available only when the model supports file parsing and understanding capability.
If the variable is of File Array type, the corresponding value should be a list whose elements contain following attributions:`]}),`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) Supported type:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," ('TXT', 'MD', 'MARKDOWN', 'PDF', 'HTML', 'XLSX', 'XLS', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB')"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," ('JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG')"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," ('MP3', 'M4A', 'WAV', 'WEBM', 'AMR')"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," ('MP4', 'MOV', 'MPEG', 'MPGA')"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (Other file types)"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) Transfer method, ",e.jsx(n.code,{children:"remote_url"})," for image URL / ",e.jsx(n.code,{children:"local_file"})," for file upload"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) Image URL (when the transfer method is ",e.jsx(n.code,{children:"remote_url"}),")"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," (string) Uploaded file ID, which must be obtained by uploading through the File Upload API in advance (when the transfer method is ",e.jsx(n.code,{children:"local_file"}),")"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"response_mode"}),` (string) Required
The mode of response return, supporting:`]}),`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," Streaming mode (recommended), implements a typewriter-like output through SSE (",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"}),")."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` Blocking mode, returns result after execution is complete. (Requests may be interrupted if the process is long)
`,e.jsx("i",{children:"Due to Cloudflare restrictions, the request will be interrupted without a return after 100 seconds."})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, used to define the identity of the end-user for retrieval and statistics.
Should be uniquely defined by the developer within the application.`]}),`
`,e.jsx("br",{}),`
`,e.jsx("i",{children:"The user identifier should be consistent with the user passed in the message sending interface. The Service API does not share conversations created by the WebApp."}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"files"})," (array[object]) Optional"]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"trace_id"}),` (string) Optional
Trace ID. Used for integration with existing business trace components to achieve end-to-end distributed tracing. If not provided, the system will automatically generate a trace_id. Supports the following three ways to pass, in order of priority:`]}),`
`,e.jsxs(n.ol,{children:[`
`,e.jsxs(n.li,{children:["Header: via HTTP Header ",e.jsx(n.code,{children:"X-Trace-Id"}),", highest priority."]}),`
`,e.jsxs(n.li,{children:["Query parameter: via URL query parameter ",e.jsx(n.code,{children:"trace_id"}),"."]}),`
`,e.jsxs(n.li,{children:["Request Body: via request body field ",e.jsx(n.code,{children:"trace_id"})," (i.e., this field)."]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.p,{children:["When ",e.jsx(n.code,{children:"response_mode"})," is ",e.jsx(n.code,{children:"blocking"}),`, return a CompletionResponse object.
When `,e.jsx(n.code,{children:"response_mode"})," is ",e.jsx(n.code,{children:"streaming"}),", return a ChunkCompletionResponse stream."]}),e.jsx(n.h3,{children:"CompletionResponse"}),e.jsxs(n.p,{children:["Returns the App result, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"application/json"}),"."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail of result",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) default 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) start time"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) end time"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"ChunkCompletionResponse"}),e.jsxs(n.p,{children:["Returns the stream chunks outputted by the App, ",e.jsx(n.code,{children:"Content-Type"})," is ",e.jsx(n.code,{children:"text/event-stream"}),`.
Each streaming chunk starts with `,e.jsx(n.code,{children:"data:"}),", separated by two newline characters ",e.jsx(n.code,{children:"\\n\\n"}),", as shown below:"]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "text_chunk", "workflow_run_id": "b85e5fc5-751b-454d-b14e-dc5f240b0a31", "task_id": "bd029338-b068-4d34-a331-fc85478922c2", "data": {"text": "\\u4e3a\\u4e86", "from_variable_selector": ["1745912968134", "text"]}}\\n\\n
`})})}),e.jsxs(n.p,{children:["The structure of the streaming chunks varies depending on the ",e.jsx(n.code,{children:"event"}),":"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," workflow starts execution",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"workflow_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," node execution started",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"node_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ID of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) type of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) name of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) Execution sequence number, used to display Tracing Node sequence"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional Prefix node ID, used for canvas display execution path"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) Contents of all preceding node variables used in the node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) timestamp of start, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: text_chunk"})," Text fragment",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"text_chunk"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) Text content"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"from_variable_selector"})," (array) Text source path, helps developers understand which variable of which node the text is generated from"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," node execution finished, success and failure are different states in the same event",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"node_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of node execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ID of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) Execution sequence number, used to display Tracing Node sequence"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional Prefix node ID, used for canvas display execution path"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) Contents of all preceding node variables used in the node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) Optional Process data of node"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) metadata",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) optional tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) optional total cost"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) optional currency, such as ",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) timestamp of start, e.g., 1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," workflow execution finished, success and failure are different states in the same event",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) fixed to ",e.jsx(n.code,{children:"workflow_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) detail",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) Unique ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) default 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) start time"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) end time"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS audio stream event, i.e., speech synthesis output. The content is an audio block in Mp3 format, encoded as a base64 string, which can be decoded directly when playing. (Only available when auto-play is enabled)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The audio block after speech synthesis is encoded as base64 text content, which can be directly base64 decoded and sent to the player when playing"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS audio stream end event, receiving this event indicates the end of audio stream return.",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, used for request tracking and the below Stop Generate API"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) Unique message ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) The end event has no audio, so this is an empty string"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) Creation timestamp, e.g., 1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," Ping event every 10s to keep the connection alive."]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", Invalid input parameters"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", App configuration unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", No available model credentials configured"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", Insufficient model call quota"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", Current model unavailable"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_not_found"}),", Specified workflow version not found"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"draft_workflow_error"}),", Cannot use draft workflow version"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_id_format_error"}),", Workflow ID format error, UUID format required"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_request_error"}),", Workflow execution failed"]}),`
`,e.jsx(n.li,{children:"500, Internal service error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/workflows/:workflow_id/run",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/{workflow_id}/run' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Example: file array as an input variable",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "inputs": {
    "{variable_name}":
    [
      {
      "transfer_method": "local_file",
      "upload_file_id": "{upload_file_id}",
      "type": "{document_type}"
      }
    ]
  }
}
`})})}),e.jsx(n.h3,{children:"Blocking Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "workflow_run_id": "djflajgkldjgd",
    "task_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "data": {
        "id": "fdlsjfjejkghjda",
        "workflow_id": "fldjaslkfjlsda",
        "status": "succeeded",
        "outputs": {
          "text": "Nice to meet you."
        },
        "error": null,
        "elapsed_time": 0.875,
        "total_tokens": 3562,
        "total_steps": 8,
        "created_at": 1705407629,
        "finished_at": 1727807631
    }
}
`})})}),e.jsx(n.h3,{children:"Streaming Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/run/:workflow_run_id",method:"GET",title:"Get Workflow Run Detail",name:"#get-workflow-run-detail"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve the current execution results of a workflow task based on the workflow execution ID."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) Workflow run ID, can be obtained from the streaming chunk return"]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID of workflow execution"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) ID of related workflow"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (json) content of input"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) content of output"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) total steps of task"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) total tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) start time"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) end time"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) total seconds to be used"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/workflows/run/:workflow_run_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/workflows/run/:workflow_run_id' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "id": "b1ad3277-089e-42c6-9dff-6820d94fbc76",
    "workflow_id": "19eff89f-ec03-4f75-b0fc-897e7effea02",
    "status": "succeeded",
    "inputs": "{\\"sys.files\\": [], \\"sys.user_id\\": \\"abc-123\\"}",
    "outputs": null,
    "error": null,
    "total_steps": 3,
    "total_tokens": 0,
    "created_at": 1705407629,
    "finished_at": 1727807631,
    "elapsed_time": 30.098514399956912
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/tasks/:task_id/stop",method:"POST",title:"Stop Generate",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Only supported in streaming mode."}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) Task ID, can be obtained from the streaming chunk return"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, used to define the identity of the end-user, must be consistent with the user passed in the message sending interface. The Service API does not share conversations created by the WebApp.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) Always returns "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/workflows/tasks/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/tasks/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{"user": "abc-123"}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"File Upload",name:"#file-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:`Upload a file for use when sending messages, enabling multimodal understanding of images and text.
Supports any formats that are supported by your workflow.
Uploaded files are for use by the current end-user only.`}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["This interface requires a ",e.jsx(n.code,{children:"multipart/form-data"})," request."]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file"}),` (File) Required
The file to be uploaded.`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
User identifier, defined by the developer's rules, must be unique within the application. The Service API does not share conversations created by the WebApp.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"After a successful upload, the server will return the file's ID and related information."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) File name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) File size (bytes)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) File extension"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) File mime-type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) End-user ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) Creation timestamp, e.g., 1705395332"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"no_file_uploaded"}),", a file must be provided"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"too_many_files"}),", currently only one file is accepted"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_preview"}),", the file does not support preview"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_estimate"}),", the file does not support estimation"]}),`
`,e.jsxs(n.li,{children:["413, ",e.jsx(n.code,{children:"file_too_large"}),", the file is too large"]}),`
`,e.jsxs(n.li,{children:["415, ",e.jsx(n.code,{children:"unsupported_file_type"}),", unsupported extension, currently only document files are accepted"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_connection_failed"}),", unable to connect to S3 service"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_permission_denied"}),", no permission to upload files to S3"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_file_too_large"}),", file exceeds S3 size limit"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"Get End User",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Retrieve an end user by ID."}),e.jsxs(n.p,{children:["This is useful when other APIs return an end-user ID (e.g. ",e.jsx(n.code,{children:"created_by"})," from File Upload)."]}),e.jsx(n.h3,{children:"Path Parameters"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) Required
End user ID.`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"Returns an EndUser object."}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) Tenant ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) App ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) End user type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) External user ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) Name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) Whether anonymous"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) Session ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 datetime"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 datetime"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"end_user_not_found"}),", end user not found"]}),`
`,e.jsx(n.li,{children:"500, internal server error"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/logs",method:"GET",title:"Get Workflow Logs",name:"#Get-Workflow-Logs"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:["Returns workflow logs, with the first page returning the latest ",e.jsx(n.code,{children:"{limit}"})," messages, i.e., in reverse order."]}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"keyword",type:"string",children:e.jsx(n.p,{children:"Keyword to search"})},"keyword"),e.jsx(d,{name:"status",type:"string",children:e.jsx(n.p,{children:"succeeded/failed/stopped"})},"status"),e.jsx(d,{name:"page",type:"int",children:e.jsx(n.p,{children:"current page, default is 1."})},"page"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"How many chat history messages to return in one request, default is 20."})},"limit"),e.jsx(d,{name:"created_by_end_user_session_id",type:"str",children:e.jsxs(n.p,{children:["Created by which endUser, for example, ",e.jsx(n.code,{children:"abc-123"}),"."]})},"created_by_end_user_session_id"),e.jsx(d,{name:"created_by_account",type:"str",children:e.jsx(n.p,{children:"Created by which email account, for example, lizb@test.com."})},"created_by_account")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"page"})," (int) Current page"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) Number of returned items, if input exceeds system limit, returns system limit amount"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total"})," (int) Number of total items"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) Whether there is a next page"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) Log list",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run"})," (object) Workflow run",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"version"})," (string) Version"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) status of execution, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional reason of error"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) total seconds to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) tokens to be used"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) default 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) start time"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) end time"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_from"})," (string) Created from"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_role"})," (string) Created by role"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_account"})," (string) Optional Created by account"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_end_user"})," (object) Created by end user",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) Type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (bool) Is anonymous"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) Session ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) create time"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/workflows/logs",targetCode:`curl -X GET '${i.appDetail.api_base_url}/workflows/logs'\\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "page": 1,
    "limit": 1,
    "total": 7,
    "has_more": true,
    "data": [
        {
            "id": "e41b93f1-7ca2-40fd-b3a8-999aeb499cc0",
            "workflow_run": {
                "id": "c0640fc8-03ef-4481-a96c-8a13b732a36e",
                "version": "2024-08-01 12:17:09.771832",
                "status": "succeeded",
                "error": null,
                "elapsed_time": 1.3588523610014818,
                "total_tokens": 0,
                "total_steps": 3,
                "created_at": 1726139643,
                "finished_at": 1726139644
            },
            "created_from": "service-api",
            "created_by_role": "end_user",
            "created_by_account": null,
            "created_by_end_user": {
                "id": "7f7d9117-dd9d-441d-8970-87e5e7e687a3",
                "type": "service_api",
                "is_anonymous": false,
                "session_id": "abc-123"
            },
            "created_at": 1726139644
        }
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"Get Application Basic Information",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get basic information about this application"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) application name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) application description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) application tags"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) application mode"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"author_name"})," (string) application author name"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "workflow",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"Get Application Parameters Information",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used at the start of entering the page to obtain information such as features, input parameter names, types, and default values."}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) User input form configuration",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) Text input control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) Paragraph text input control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) Dropdown control",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) Variable display label name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) Variable ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) Whether it is required"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) Default value"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) Option values"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) File upload configuration",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) Document settings
Currently only supports document types: `,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Document number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) Image settings
Currently only supports image types: `,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Image number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) Audio settings
Currently only supports audio types: `,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Audio number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) Video settings
Currently only supports video types: `,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),".",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Video number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) Custom settings",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) Whether it is enabled"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) Custom number limit, default is 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) List of transfer methods: ",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),". Must choose one."]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) System parameters",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) Document upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) Image file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) Audio file upload size limit (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) Video file upload size limit (MB)"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "user_input_form": [
      {
          "paragraph": {
              "label": "Query",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"Get Application WebApp Settings",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"Used to get the WebApp settings of the application."}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp name"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) Icon type, ",e.jsx(n.code,{children:"emoji"})," - emoji, ",e.jsx(n.code,{children:"image"})," - picture"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) Icon. If it's ",e.jsx(n.code,{children:"emoji"})," type, it's an emoji symbol; if it's ",e.jsx(n.code,{children:"image"})," type, it's an image URL."]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) Background color in hex format"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) Icon URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) Description"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) Copyright information"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) Privacy policy link"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) Custom disclaimer"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) Default language"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) Whether to show workflow details"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{})]})}function Qn(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(ve,{...i})}):ve(i)}function ke(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"ワークフローアプリ API"}),`
`,e.jsx(n.p,{children:"ワークフローアプリケーションは、セッションをサポートせず、翻訳、記事作成、要約 AI などに最適です。"}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"ベース URL"}),e.jsx(s,{title:"コード",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"認証"}),e.jsxs(n.p,{children:["サービス API は ",e.jsx(n.code,{children:"API-Key"}),` 認証を使用します。
`,e.jsx("i",{children:e.jsx(n.strong,{children:"API キーの漏洩を防ぐため、API キーはクライアント側で共有または保存せず、サーバー側で保存することを強くお勧めします。"})})]}),e.jsxs(n.p,{children:["すべての API リクエストにおいて、以下のように ",e.jsx(n.code,{children:"Authorization"}),"HTTP ヘッダーに API キーを含めてください："]}),e.jsx(s,{title:"コード",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/run",method:"POST",title:"ワークフローを実行",name:"#Execute-Workflow"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ワークフローを実行します。公開されたワークフローがないと実行できません。"}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"inputs"}),` (object) 必須
アプリで定義されたさまざまな変数値の入力を許可します。
`,e.jsx(n.code,{children:"inputs"}),`パラメータには複数のキー/値ペアが含まれ、各キーは特定の変数に対応し、各値はその変数の特定の値です。
ワークフローアプリケーションは少なくとも1つのキー/値ペアの入力を必要とします。値はファイルリストである場合もあります。
ファイルリストは、テキスト理解と質問への回答を組み合わせたファイルの入力に適しています。モデルがファイルの解析と理解機能をサポートしている場合にのみ使用できます。`]}),`
`,e.jsx(n.p,{children:"変数がファイルリストの場合、リストの各要素は以下の属性を持つ必要があります。"}),`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) サポートされるタイプ：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," サポートされるタイプには以下が含まれます：'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," サポートされるタイプには以下が含まれます：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," サポートされるタイプには以下が含まれます：'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," サポートされるタイプには以下が含まれます：'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," サポートされるタイプには以下が含まれます：その他のファイルタイプ"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) 転送方法:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": ファイルのURL。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": ファイルをアップロード。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," ファイルのURL。（転送方法が ",e.jsx(n.code,{children:"remote_url"})," の場合のみ）。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," アップロードされたファイルID。（転送方法が ",e.jsx(n.code,{children:"local_file"})," の場合のみ）。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"response_mode"}),` (string) 必須
応答の返却モードを指定します。サポートされているモード：`]}),`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," ストリーミングモード（推奨）、SSE（",e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"}),"）を通じてタイプライターのような出力を実装します。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` ブロッキングモード、実行完了後に結果を返します。（プロセスが長い場合、リクエストが中断される可能性があります）
`,e.jsx("i",{children:"Cloudflare の制限により、100 秒後に応答がない場合、リクエストは中断されます。"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子、エンドユーザーのアイデンティティを定義するために使用されます。
アプリケーション内で開発者によって一意に定義される必要があります。`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"files"})," (array[object]) オプション"]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"trace_id"}),` (string) オプション
トレースID。既存の業務システムのトレースコンポーネントと連携し、エンドツーエンドの分散トレーシングを実現するために使用します。指定がない場合、システムが自動的に trace_id を生成します。以下の3つの方法で渡すことができ、優先順位は次のとおりです：`]}),`
`,e.jsxs(n.ol,{children:[`
`,e.jsxs(n.li,{children:["Header：HTTPヘッダー ",e.jsx(n.code,{children:"X-Trace-Id"})," で渡す（最優先）。"]}),`
`,e.jsxs(n.li,{children:["クエリパラメータ：URLクエリパラメータ ",e.jsx(n.code,{children:"trace_id"})," で渡す。"]}),`
`,e.jsxs(n.li,{children:["リクエストボディ：リクエストボディの ",e.jsx(n.code,{children:"trace_id"})," フィールドで渡す（本フィールド）。"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.p,{children:[e.jsx(n.code,{children:"response_mode"}),"が",e.jsx(n.code,{children:"blocking"}),`の場合、CompletionResponse オブジェクトを返します。
`,e.jsx(n.code,{children:"response_mode"}),"が",e.jsx(n.code,{children:"streaming"}),"の場合、ChunkCompletionResponse ストリームを返します。"]}),e.jsx(n.h3,{children:"CompletionResponse"}),e.jsxs(n.p,{children:["アプリの結果を返します。",e.jsx(n.code,{children:"Content-Type"}),"は",e.jsx(n.code,{children:"application/json"}),"です。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、リクエスト追跡と以下の Stop Generate API に使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 結果の詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連するワークフローの ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行のステータス、",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) オプションの出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプションのエラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) オプションの使用時間（秒）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) オプションの使用トークン数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) デフォルト 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 終了時間"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"ChunkCompletionResponse"}),e.jsxs(n.p,{children:["アプリによって出力されたストリームチャンクを返します。",e.jsx(n.code,{children:"Content-Type"}),"は",e.jsx(n.code,{children:"text/event-stream"}),`です。
各ストリーミングチャンクは`,e.jsx(n.code,{children:"data:"}),"で始まり、2 つの改行文字",e.jsx(n.code,{children:"\\n\\n"}),"で区切られます。以下のように表示されます："]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "text_chunk", "workflow_run_id": "b85e5fc5-751b-454d-b14e-dc5f240b0a31", "task_id": "bd029338-b068-4d34-a331-fc85478922c2", "data": {"text": "\\u4e3a\\u4e86", "from_variable_selector": ["1745912968134", "text"]}}\\n\\n
`})})}),e.jsxs(n.p,{children:["ストリーミングチャンクの構造は",e.jsx(n.code,{children:"event"}),"に応じて異なります："]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," ワークフローが実行を開始",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、リクエスト追跡と以下の Stop Generate API に使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"workflow_started"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連するワークフローの ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," ノード実行開始",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、リクエスト追跡と以下の Stop Generate API に使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"node_started"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ノードの ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) ノードのタイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) ノードの名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 実行シーケンス番号、トレースノードシーケンスを表示するために使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) オプションのプレフィックスノード ID、キャンバス表示実行パスに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ノードで使用されるすべての前のノード変数の内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始のタイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: text_chunk"})," テキストフラグメント",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、リクエスト追跡と以下の Stop Generate API に使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"text_chunk"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) テキスト内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"from_variable_selector"})," (array) テキスト生成元パス（開発者がどのノードのどの変数から生成されたかを理解するための情報）"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," ノード実行終了、同じイベントで異なる状態で成功または失敗",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、リクエスト追跡と以下の Stop Generate API に使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"node_finished"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ノードの ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) ノードのタイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) ノードの名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 実行シーケンス番号、トレースノードシーケンスを表示するために使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) オプションのプレフィックスノード ID、キャンバス表示実行パスに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ノードで使用されるすべての前のノード変数の内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) オプションのノードプロセスデータ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) オプションの出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行のステータス、",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプションのエラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) オプションの使用時間（秒）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) オプションの使用トークン数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) オプションの総コスト"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) オプション 例：",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始のタイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," ワークフロー実行終了、同じイベントで異なる状態で成功または失敗",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、リクエスト追跡と以下の Stop Generate API に使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行の一意の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"workflow_finished"}),"に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連するワークフローの ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行のステータス、",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) オプションの出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプションのエラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) オプションの使用時間（秒）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) オプションの使用トークン数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) デフォルト 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 終了時間"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS オーディオストリームイベント、つまり音声合成出力。内容は Mp3 形式のオーディオブロックで、base64 文字列としてエンコードされています。再生時には、base64 をデコードしてプレーヤーに入力するだけです。（このメッセージは自動再生が有効な場合にのみ利用可能）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージ ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 音声合成後のオーディオ、base64 テキストコンテンツとしてエンコードされており、再生時には base64 をデコードしてプレーヤーに入力するだけです"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS オーディオストリーム終了イベント。このイベントを受信すると、オーディオストリームの終了を示します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 一意のメッセージ ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 終了イベントにはオーディオがないため、これは空の文字列です"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 接続を維持するために 10 秒ごとに送信される Ping イベント。"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"invalid_param"}),", 異常なパラメータ入力"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"app_unavailable"}),", アプリの設定が利用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_not_initialize"}),", 利用可能なモデル資格情報の設定がありません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"provider_quota_exceeded"}),", モデル呼び出しのクォータが不足しています"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"model_currently_not_support"}),", 現在のモデルは利用できません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"workflow_request_error"}),", ワークフロー実行に失敗しました"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/workflows/run",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/run' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"ファイル変数の例",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "inputs": {
    "{variable_name}":
    [
      {
      "transfer_method": "local_file",
      "upload_file_id": "{upload_file_id}",
      "type": "{document_type}"
      }
    ]
  }
}
`})})}),e.jsx(n.h3,{children:"ブロッキングモード"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "workflow_run_id": "djflajgkldjgd",
    "task_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "data": {
        "id": "fdlsjfjejkghjda",
        "workflow_id": "fldjaslkfjlsda",
        "status": "succeeded",
        "outputs": {
          "text": "Nice to meet you."
        },
        "error": null,
        "elapsed_time": 0.875,
        "total_tokens": 3562,
        "total_steps": 8,
        "created_at": 1705407629,
        "finished_at": 1727807631
    }
}
`})})}),e.jsx(n.h3,{children:"ストリーミングモード"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})}),e.jsx(s,{title:"ファイルアップロードのサンプルコード",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`import requests
import json

def upload_file(file_path, user):
    upload_url = "https://api.gofy.ai/v1/files/upload"
    headers = {
        "Authorization": "Bearer app-xxxxxxxx",
    }

    try:
        print("ファイルをアップロードしています...")
        with open(file_path, 'rb') as file:
            files = {
                'file': (file_path, file, 'text/plain')  # ファイルが適切な MIME タイプでアップロードされていることを確認してください
            }
            data = {
                "user": user,
                "type": "TXT"  # ファイルタイプをTXTに設定します
            }

            response = requests.post(upload_url, headers=headers, files=files, data=data)
            if response.status_code == 201:  # 201 は作成が成功したことを意味します
                print("ファイルが正常にアップロードされました")
                return response.json().get("id")  # アップロードされたファイルIDを取得する
            else:
                print(f"ファイルのアップロードに失敗しました。ステータス コード: {response.status_code}")
                return None
    except Exception as e:
        print(f"エラーが発生しました: {str(e)}")
        return None

def run_workflow(file_id, user, response_mode="blocking"):
    workflow_url = "https://api.gofy.ai/v1/workflows/run"
    headers = {
        "Authorization": "Bearer app-xxxxxxxxx",
        "Content-Type": "application/json"
    }

    data = {
        "inputs": {
            "orig_mail": [{
                "transfer_method": "local_file",
                "upload_file_id": file_id,
                "type": "document"
            }]
        },
        "response_mode": response_mode,
        "user": user
    }

    try:
        print("ワークフローを実行...")
        response = requests.post(workflow_url, headers=headers, json=data)
        if response.status_code == 200:
            print("ワークフローが正常に実行されました")
            return response.json()
        else:
            print(f"ワークフローの実行がステータス コードで失敗しました: {response.status_code}")
            return {"status": "error", "message": f"Failed to execute workflow, status code: {response.status_code}"}
    except Exception as e:
        print(f"エラーが発生しました: {str(e)}")
        return {"status": "error", "message": str(e)}

# 使用例
file_path = "{your_file_path}"
user = "gofyuser"

# ファイルをアップロードする
file_id = upload_file(file_path, user)
if file_id:
    # ファイルは正常にアップロードされました。ワークフローの実行を続行します
    result = run_workflow(file_id, user)
    print(result)
else:
    print("ファイルのアップロードに失敗し、ワークフローを実行できません")
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/:workflow_id/run",method:"POST",title:"特定バージョンのワークフローを実行",name:"#Execute-Specific-Workflow"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"パスパラメータでワークフローIDを指定して、特定バージョンのワークフローを実行します。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 必須 特定バージョンのワークフローを指定するためのワークフローID"]}),`
`,e.jsx(n.p,{children:"取得方法：バージョン履歴インターフェースで、各バージョンエントリの右側にあるコピーアイコンをクリックすると、完全なワークフローIDをコピーできます。"}),`
`]}),`
`]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"}),` (object) 必須
App で定義された各変数値を入力できます。
inputs パラメータには複数のキー/値ペアが含まれており、各キーは特定の変数に対応し、各値はその変数の具体的な値です。変数はファイルリスト型にすることができます。
ファイルリスト型変数は、ファイルをテキスト理解と組み合わせて質問に答えるために入力するのに適しており、モデルがファイル解析機能をサポートしている場合のみ使用できます。変数がファイルリスト型の場合、その変数に対応する値はリスト形式である必要があり、各要素には以下の内容が含まれます：`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) サポートされるタイプ：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," 具体的なタイプには以下が含まれます：'TXT', 'MD', 'MARKDOWN', 'PDF', 'HTML', 'XLSX', 'XLS', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," 具体的なタイプには以下が含まれます：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," 具体的なタイプには以下が含まれます：'MP3', 'M4A', 'WAV', 'WEBM', 'AMR'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," 具体的なタイプには以下が含まれます：'MP4', 'MOV', 'MPEG', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," 具体的なタイプには以下が含まれます：その他のファイルタイプ"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) 転送方法、",e.jsx(n.code,{children:"remote_url"})," 画像URL / ",e.jsx(n.code,{children:"local_file"})," ファイルアップロード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) 画像URL（転送方法が ",e.jsx(n.code,{children:"remote_url"})," の場合のみ）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," (string) アップロードされたファイルID（転送方法が ",e.jsx(n.code,{children:"local_file"})," の場合のみ）"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"response_mode"}),` (string) 必須
応答返却モード、以下をサポート：`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," ストリーミングモード（推奨）。SSE（",e.jsx(n.strong,{children:e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"})}),"）をベースにタイプライター風の出力を実現。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` ブロッキングモード、実行完了後に結果を返却。（プロセスが長い場合、リクエストが中断される可能性があります）。
`,e.jsx("i",{children:"Cloudflare の制限により、100秒後に応答がない場合、リクエストは中断されます。"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子、エンドユーザーのアイデンティティを定義し、検索・統計を容易にするために使用されます。
開発者が定義するルールで、アプリケーション内でユーザー識別子が一意である必要があります。API は WebApp で作成されたセッションにアクセスできません。`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"files"})," (array[object]) オプション"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"trace_id"}),` (string) オプション
トレースID。既存のビジネスシステムのトレースコンポーネントと統合して、エンドツーエンドの分散トレーシングを実現するために使用されます。指定されていない場合、システムは自動的に `,e.jsx(n.code,{children:"trace_id"})," を生成します。以下の3つの方法で渡すことができ、優先順位は以下の通りです：",`
`,e.jsxs(n.ol,{children:[`
`,e.jsxs(n.li,{children:["ヘッダー：HTTP ヘッダー ",e.jsx(n.code,{children:"X-Trace-Id"})," で渡すことを推奨、最高優先度。"]}),`
`,e.jsxs(n.li,{children:["クエリパラメータ：URL クエリパラメータ ",e.jsx(n.code,{children:"trace_id"})," で渡す。"]}),`
`,e.jsxs(n.li,{children:["リクエストボディ：リクエストボディフィールド ",e.jsx(n.code,{children:"trace_id"})," で渡す（つまり、このフィールド）。"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.p,{children:[e.jsx(n.code,{children:"response_mode"})," が ",e.jsx(n.code,{children:"blocking"}),` の場合、CompletionResponse オブジェクトを返します。
`,e.jsx(n.code,{children:"response_mode"})," が ",e.jsx(n.code,{children:"streaming"})," の場合、ChunkCompletionResponse オブジェクトのストリーミングシーケンスを返します。"]}),e.jsx(n.h3,{children:"CompletionResponse"}),e.jsxs(n.p,{children:["完全な App 結果を返し、",e.jsx(n.code,{children:"Content-Type"})," は ",e.jsx(n.code,{children:"application/json"})," です。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連するワークフローID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行ステータス、",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) オプション 出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプション エラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) オプション 使用時間(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) オプション 使用されるトークンの総数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 総ステップ数（冗長）、デフォルト 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 終了時間"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"ChunkCompletionResponse"}),e.jsxs(n.p,{children:["App の出力ストリーミングチャンクを返し、",e.jsx(n.code,{children:"Content-Type"})," は ",e.jsx(n.code,{children:"text/event-stream"}),` です。
各ストリーミングチャンクは `,e.jsx(n.code,{children:"data:"})," で始まり、チャンク間は ",e.jsx(n.code,{children:"\\n\\n"})," つまり2つの改行文字で区切られます。以下のようになります："]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "text_chunk", "workflow_run_id": "b85e5fc5-751b-454d-b14e-dc5f240b0a31", "task_id": "bd029338-b068-4d34-a331-fc85478922c2", "data": {"text": "\\u4e3a\\u4e86", "from_variable_selector": ["1745912968134", "text"]}}\\n\\n
`})})}),e.jsxs(n.p,{children:["ストリーミングチャンクは ",e.jsx(n.code,{children:"event"})," によって構造が異なり、以下のタイプが含まれます："]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," ワークフロー実行開始",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"workflow_started"})," に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連するワークフローID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," ノード実行開始",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"node_started"})," に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ノードID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) ノードタイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) ノード名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 実行シーケンス番号、Tracing Node シーケンスの表示に使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) 前置ノードID、キャンバス表示実行パスに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ノードで使用されるすべての前置ノード変数の内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: text_chunk"})," テキストフラグメント",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"text_chunk"})," に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) テキスト内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"from_variable_selector"})," (array) テキストソースパス、開発者がテキストがどのノードのどの変数から生成されたかを理解するのに役立ちます"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," ノード実行終了、成功と失敗は同じイベント内の異なる状態",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"node_finished"})," に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ノード実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) ノードID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 実行シーケンス番号、Tracing Node シーケンスの表示に使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) オプション 前置ノードID、キャンバス表示実行パスに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) ノードで使用されるすべての前置ノード変数の内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) オプション ノードプロセスデータ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) オプション 出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行ステータス ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプション エラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) オプション 使用時間(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) メタデータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) オプション 使用されるトークンの総数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) オプション 総費用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) オプション 通貨、例：",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," ワークフロー実行終了、成功と失敗は同じイベント内の異なる状態",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) ",e.jsx(n.code,{children:"workflow_finished"})," に固定"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 詳細内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連するワークフローID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行ステータス ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) オプション 出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプション エラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) オプション 使用時間(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) オプション 使用されるトークンの総数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 総ステップ数（冗長）、デフォルト 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 終了時間"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS オーディオストリームイベント、つまり：音声合成出力。内容はMp3形式のオーディオブロックで、base64エンコードされた文字列として、再生時に直接デコードできます。（自動再生が有効な場合のみこのメッセージがあります）",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) メッセージ一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 音声合成後のオーディオブロックはbase64エンコードされたテキスト内容として、再生時に直接base64デコードしてプレーヤーに送信できます"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS オーディオストリーム終了イベント、このイベントを受信すると、オーディオストリームの返却が終了したことを示します。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスクID、リクエスト追跡と以下の停止応答インターフェースに使用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) メッセージ一意ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 終了イベントにはオーディオがないため、ここは空文字列です"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 作成タイムスタンプ、例：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 10秒ごとのpingイベント、接続を維持します。"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"invalid_param"}),"，入力パラメータ異常"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"app_unavailable"}),"，App 設定が利用できません"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_not_initialize"}),"，利用可能なモデル認証情報設定がありません"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_quota_exceeded"}),"，モデル呼び出しクォータが不足しています"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"model_currently_not_support"}),"，現在のモデルが利用できません"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_not_found"}),"，指定されたワークフローバージョンが見つかりません"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"draft_workflow_error"}),"，ドラフトワークフローバージョンを使用できません"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_id_format_error"}),"，ワークフローID形式エラー、UUID形式が必要です"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_request_error"}),"，ワークフロー実行に失敗しました"]}),`
`,e.jsx(n.li,{children:"500，サービス内部異常"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"POST",label:"/workflows/:workflow_id/run",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/{workflow_id}/run' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"例：入力変数としてのファイル配列",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "inputs": {
    "{variable_name}":
    [
      {
      "transfer_method": "local_file",
      "upload_file_id": "{upload_file_id}",
      "type": "{document_type}"
      }
    ]
  }
}
`})})}),e.jsx(n.h3,{children:"ブロッキングモード"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "workflow_run_id": "djflajgkldjgd",
    "task_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "data": {
        "id": "fdlsjfjejkghjda",
        "workflow_id": "fldjaslkfjlsda",
        "status": "succeeded",
        "outputs": {
          "text": "Nice to meet you."
        },
        "error": null,
        "elapsed_time": 0.875,
        "total_tokens": 3562,
        "total_steps": 8,
        "created_at": 1705407629,
        "finished_at": 1727807631
    }
}
`})})}),e.jsx(n.h3,{children:"ストリーミングモード"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/run/:workflow_run_id",method:"GET",title:"ワークフロー実行詳細を取得",name:"#get-workflow-run-detail"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ワークフロー実行 ID に基づいて、ワークフロータスクの現在の実行結果を取得します。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) ワークフロー実行ID、ストリーミングチャンクの返り値から取得可能"]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ワークフロー実行の ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 関連するワークフローの ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行のステータス、",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (json) 入力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) 出力内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) エラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) タスクの総ステップ数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) 使用されるトークンの総数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 終了時間"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) 使用される総秒数"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"GET",label:"/workflows/run/:workflow_run_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/workflows/run/:workflow_run_id' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "id": "b1ad3277-089e-42c6-9dff-6820d94fbc76",
    "workflow_id": "19eff89f-ec03-4f75-b0fc-897e7effea02",
    "status": "succeeded",
    "inputs": "{\\"sys.files\\": [], \\"sys.user_id\\": \\"abc-123\\"}",
    "outputs": null,
    "error": null,
    "total_steps": 3,
    "total_tokens": 0,
    "created_at": 1705407629,
    "finished_at": 1727807631,
    "elapsed_time": 30.098514399956912
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/tasks/:task_id/stop",method:"POST",title:"生成を停止",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ストリーミングモードでのみサポートされています。"}),e.jsx(n.h3,{children:"パス"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) タスク ID、ストリーミングチャンクの返り値から取得可能"]}),`
`]}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子、エンドユーザーのアイデンティティを定義するために使用され、送信メッセージインターフェースで渡されたユーザーと一致している必要があります。サービス API は WebApp によって作成された会話を共有しません。`]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 常に"success"を返します']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"POST",label:"/workflows/tasks/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/tasks/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{"user": "abc-123"}'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"ファイルアップロード",name:"#file-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:`メッセージ送信時に使用するためのファイルをアップロードし、画像とテキストのマルチモーダル理解を可能にします。
ワークフローでサポートされている任意の形式をサポートします。
アップロードされたファイルは、現在のエンドユーザーのみが使用できます。`}),e.jsx(n.h3,{children:"リクエストボディ"}),e.jsxs(n.p,{children:["このインターフェースは",e.jsx(n.code,{children:"multipart/form-data"}),"リクエストを必要とします。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file"}),` (File) 必須
アップロードするファイル。`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) 必須
ユーザー識別子、開発者のルールで定義され、アプリケーション内で一意でなければなりません。サービス API は WebApp によって作成された会話を共有しません。`]}),`
`]}),e.jsx(n.h3,{children:"応答"}),e.jsx(n.p,{children:"アップロードが成功すると、サーバーはファイルの ID と関連情報を返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) ファイル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) ファイルサイズ（バイト）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) ファイル拡張子"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) ファイルの MIME タイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) エンドユーザーID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成タイムスタンプ、例：1705395332"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"no_file_uploaded"}),", ファイルが提供されていません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"too_many_files"}),", 現在は 1 つのファイルのみ受け付けています"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_preview"}),", ファイルはプレビューをサポートしていません"]}),`
`,e.jsxs(n.li,{children:["400, ",e.jsx(n.code,{children:"unsupported_estimate"}),", ファイルは推定をサポートしていません"]}),`
`,e.jsxs(n.li,{children:["413, ",e.jsx(n.code,{children:"file_too_large"}),", ファイルが大きすぎます"]}),`
`,e.jsxs(n.li,{children:["415, ",e.jsx(n.code,{children:"unsupported_file_type"}),", サポートされていない拡張子、現在はドキュメントファイルのみ受け付けています"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_connection_failed"}),", S3 サービスに接続できません"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_permission_denied"}),", S3 にファイルをアップロードする権限がありません"]}),`
`,e.jsxs(n.li,{children:["503, ",e.jsx(n.code,{children:"s3_file_too_large"}),", ファイルが S3 のサイズ制限を超えています"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"リクエスト",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"エンドユーザーを取得",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"エンドユーザー ID からエンドユーザー情報を取得します。"}),e.jsxs(n.p,{children:["他の API がエンドユーザー ID（例：ファイルアップロードの ",e.jsx(n.code,{children:"created_by"}),"）を返す場合に利用できます。"]}),e.jsx(n.h3,{children:"パスパラメータ"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) 必須
エンドユーザー ID。`]}),`
`]}),e.jsx(n.h3,{children:"レスポンス"}),e.jsx(n.p,{children:"EndUser オブジェクトを返します。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) テナント ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) アプリ ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) エンドユーザー種別"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) 外部ユーザー ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) 匿名ユーザーかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) セッション ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 日時"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 日時"]}),`
`]}),e.jsx(n.h3,{children:"エラー"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404, ",e.jsx(n.code,{children:"end_user_not_found"}),", エンドユーザーが見つかりません"]}),`
`,e.jsx(n.li,{children:"500, 内部サーバーエラー"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"リクエスト例"}),e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"レスポンス例"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/logs",method:"GET",title:"ワークフローログを取得",name:"#Get-Workflow-Logs"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:["ワークフローログを返します。最初のページは最新の",e.jsx(n.code,{children:"{limit}"}),"メッセージを返します。つまり、逆順です。"]}),e.jsx(n.h3,{children:"クエリ"}),e.jsxs(t,{children:[e.jsx(d,{name:"keyword",type:"string",children:e.jsx(n.p,{children:"検索するキーワード"})},"keyword"),e.jsx(d,{name:"status",type:"string",children:e.jsx(n.p,{children:"succeeded/failed/stopped"})},"status"),e.jsx(d,{name:"page",type:"int",children:e.jsx(n.p,{children:"現在のページ、デフォルトは1。"})},"page"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"1回のリクエストで返すチャット履歴メッセージの数、デフォルトは20。"})},"limit"),e.jsx(d,{name:"created_by_end_user_session_id",type:"str",children:e.jsxs(n.p,{children:["どのendUserによって作成されたか、例えば、",e.jsx(n.code,{children:"abc-123"}),"。"]})},"created_by_end_user_session_id"),e.jsx(d,{name:"created_by_account",type:"str",children:e.jsx(n.p,{children:"どのメールアカウントによって作成されたか、例えば、lizb@test.com。"})},"created_by_account")]}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"page"})," (int) 現在のページ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 返されたアイテムの数、入力がシステム制限を超える場合、システム制限量を返します"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total"})," (int) 合計アイテム数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) 次のページがあるかどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) ログリスト",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run"})," (object) ワークフロー実行",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"version"})," (string) バージョン"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 実行のステータス、",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) オプションのエラー理由"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) 使用される総秒数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) 使用されるトークン数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) デフォルト 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 開始時間"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 終了時間"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_from"})," (string) 作成元"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_role"})," (string) 作成者の役割"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_account"})," (string) オプションの作成者アカウント"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_end_user"})," (object) エンドユーザーによって作成",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) タイプ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (bool) 匿名かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) セッション ID"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 作成時間"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/workflows/logs",targetCode:`curl -X GET '${i.appDetail.api_base_url}/workflows/logs'\\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"応答例"}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "page": 1,
    "limit": 1,
    "total": 7,
    "has_more": true,
    "data": [
        {
            "id": "e41b93f1-7ca2-40fd-b3a8-999aeb499cc0",
            "workflow_run": {
                "id": "c0640fc8-03ef-4481-a96c-8a13b732a36e",
                "version": "2024-08-01 12:17:09.771832",
                "status": "succeeded",
                "error": null,
                "elapsed_time": 1.3588523610014818,
                "total_tokens": 0,
                "total_steps": 3,
                "created_at": 1726139643,
                "finished_at": 1726139644
            },
            "created_from": "service-api",
            "created_by_role": "end_user",
            "created_by_account": null,
            "created_by_end_user": {
                "id": "7f7d9117-dd9d-441d-8970-87e5e7e687a3",
                "type": "service_api",
                "is_anonymous": false,
                "session_id": "abc-123"
            },
            "created_at": 1726139644
        }
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"アプリケーションの基本情報を取得",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"このアプリケーションの基本情報を取得するために使用されます"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) アプリケーションの名前"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) アプリケーションの説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) アプリケーションのタグ"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) アプリケーションのモード"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"author_name"})," (string) 作者の名前"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "workflow",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"アプリケーションのパラメータ情報を取得",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"ページに入る際に、機能、入力パラメータ名、タイプ、デフォルト値などの情報を取得するために使用されます。"}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) ユーザー入力フォームの設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) テキスト入力コントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) 段落テキスト入力コントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) ドロップダウンコントロール",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 変数表示ラベル名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 変数ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 必須かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) デフォルト値"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) オプション値"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) ファイルアップロード設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) ドキュメント設定
現在サポートされているドキュメントタイプ：`,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) ドキュメント数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) 画像設定
現在サポートされている画像タイプ：`,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 画像数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) オーディオ設定
現在サポートされているオーディオタイプ：`,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) オーディオ数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) ビデオ設定
現在サポートされているビデオタイプ：`,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) ビデオ数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) カスタム設定",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 有効かどうか"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) カスタム数の上限。デフォルトは 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 転送方法リスト：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"。いずれかを選択する必要があります。"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) システムパラメータ",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) ドキュメントアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) 画像ファイルアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) オーディオファイルアップロードサイズ制限（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) ビデオファイルアップロードサイズ制限（MB）"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"リクエスト",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"応答",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "user_input_form": [
      {
          "paragraph": {
              "label": "Query",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.p,{children:"———"}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"アプリのWebApp設定を取得",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"アプリの WebApp 設定を取得するために使用します。"}),e.jsx(n.h3,{children:"応答"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp 名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) アイコンタイプ、",e.jsx(n.code,{children:"emoji"}),"-絵文字、",e.jsx(n.code,{children:"image"}),"-画像"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) アイコン。",e.jsx(n.code,{children:"emoji"}),"タイプの場合は絵文字、",e.jsx(n.code,{children:"image"}),"タイプの場合は画像 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) 16 進数形式の背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) アイコンの URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 説明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) 著作権情報"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) プライバシーポリシーのリンク"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) カスタム免責事項"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) デフォルト言語"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) ワークフローの詳細を表示するかどうか"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{})]})}function Jn(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(ke,{...i})}):ke(i)}function we(i){const n={a:"a",code:"code",h1:"h1",h3:"h3",hr:"hr",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...i.components};return e.jsxs(e.Fragment,{children:[e.jsx(n.h1,{children:"Workflow 应用 API"}),`
`,e.jsx(n.p,{children:"Workflow 应用无会话支持，适合用于翻译/文章写作/总结 AI 等等。"}),`
`,e.jsxs("div",{children:[e.jsx(n.h3,{children:"Base URL"}),e.jsx(s,{title:"Code",targetCode:i.appDetail.api_base_url}),e.jsx(n.h3,{children:"Authentication"}),e.jsxs(n.p,{children:["Service API 使用 ",e.jsx(n.code,{children:"API-Key"}),` 进行鉴权。
`,e.jsx("i",{children:e.jsxs(n.strong,{children:["强烈建议开发者把 ",e.jsx(n.code,{children:"API-Key"})," 放在后端存储，而非分享或者放在客户端存储，以免 ",e.jsx(n.code,{children:"API-Key"})," 泄露，导致财产损失。"]})}),`
所有 API 请求都应在 `,e.jsx(n.strong,{children:e.jsx(n.code,{children:"Authorization"})})," HTTP Header 中包含您的 ",e.jsx(n.code,{children:"API-Key"}),"，如下所示："]}),e.jsx(s,{title:"Code",targetCode:"Authorization: Bearer {API_KEY}"})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/run",method:"POST",title:"执行 Workflow",name:"#Execute-Workflow"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"执行 workflow，没有已发布的 workflow，不可执行。"}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"}),` (object) Required
允许传入 App 定义的各变量值。
inputs 参数包含了多组键值对（Key/Value pairs），每组的键对应一个特定变量，每组的值则是该变量的具体值。变量可以是文件列表类型。
文件列表类型变量适用于传入文件结合文本理解并回答问题，仅当模型支持该类型文件解析能力时可用。如果该变量是文件列表类型，该变量对应的值应是列表格式，其中每个元素应包含以下内容：`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 支持类型：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," 具体类型包含：'TXT', 'MD', 'MARKDOWN', 'PDF', 'HTML', 'XLSX', 'XLS', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," 具体类型包含：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," 具体类型包含：'MP3', 'M4A', 'WAV', 'WEBM', 'AMR'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," 具体类型包含：'MP4', 'MOV', 'MPEG', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," 具体类型包含：其他文件类型"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string) 传递方式，",e.jsx(n.code,{children:"remote_url"})," 图片地址 / ",e.jsx(n.code,{children:"local_file"})," 上传文件"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," (string) 图片地址（仅当传递方式为 ",e.jsx(n.code,{children:"remote_url"})," 时）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," (string)  上传文件 ID（仅当传递方式为 ",e.jsx(n.code,{children:"local_file"})," 时）"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"response_mode"}),` (string) Required
返回响应模式，支持：`,`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," 流式模式（推荐）。基于 SSE（",e.jsx(n.strong,{children:e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"})}),"）实现类似打字机输出方式的流式返回。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` 阻塞模式，等待执行完毕后返回结果。（请求若流程较长可能会被中断）。
`,e.jsx("i",{children:"由于 Cloudflare 限制，请求会在 100 秒超时无返回后中断。"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
用户标识，用于定义终端用户的身份，方便检索、统计。
由开发者定义规则，需保证用户标识在应用内唯一。API 无法访问 WebApp 创建的会话。`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"files"})," (array[object]) 可选"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"trace_id"}),` (string) Optional
链路追踪ID。适用于与业务系统已有的trace组件打通，实现端到端分布式追踪等场景。如果未指定，系统将自动生成 `,e.jsx(n.code,{children:"trace_id"}),"。支持以下三种方式传递，具体优先级依次为：",`
`,e.jsxs(n.ol,{children:[`
`,e.jsxs(n.li,{children:["Header：推荐通过 HTTP Header ",e.jsx(n.code,{children:"X-Trace-Id"})," 传递，优先级最高。"]}),`
`,e.jsxs(n.li,{children:["Query 参数：通过 URL 查询参数 ",e.jsx(n.code,{children:"trace_id"})," 传递。"]}),`
`,e.jsxs(n.li,{children:["Request Body：通过请求体字段 ",e.jsx(n.code,{children:"trace_id"})," 传递（即本字段）。"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.p,{children:["当 ",e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"blocking"}),` 时，返回 CompletionResponse object。
当 `,e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"streaming"}),"时，返回 ChunkCompletionResponse object 流式序列。"]}),e.jsx(n.h3,{children:"CompletionResponse"}),e.jsxs(n.p,{children:["返回完整的 App 结果，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"application/json"})," 。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 执行状态, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional 输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional 耗时(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional 总使用 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 总步数（冗余），默认 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 结束时间"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"ChunkCompletionResponse"}),e.jsxs(n.p,{children:["返回 App 输出的流式块，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"text/event-stream"}),`。
每个流式块均为 data: 开头，块之间以 `,e.jsx(n.code,{children:"\\n\\n"})," 即两个换行符分隔，如下所示："]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "text_chunk", "workflow_run_id": "b85e5fc5-751b-454d-b14e-dc5f240b0a31", "task_id": "bd029338-b068-4d34-a331-fc85478922c2", "data": {"text": "\\u4e3a\\u4e86", "from_variable_selector": ["1745912968134", "text"]}}\\n\\n
`})})}),e.jsxs(n.p,{children:["流式块中根据 ",e.jsx(n.code,{children:"event"})," 不同，结构也不同，包含以下类型："]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," workflow 开始执行",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"workflow_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," node 开始执行",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"node_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) 节点 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) 节点类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) 节点名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 执行序号，用于展示 Tracing Node 顺序"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) 前置节点 ID，用于画布展示执行路径"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 节点中所有使用到的前置节点变量内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: text_chunk"})," 文本片段",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"text_chunk"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) 文本内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"from_variable_selector"})," (array) 文本来源路径，帮助开发者了解文本是由哪个节点的哪个变量生成的"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," node 执行结束，成功失败同一事件中不同状态",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"node_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) node 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) 节点 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 执行序号，用于展示 Tracing Node 顺序"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional 前置节点 ID，用于画布展示执行路径"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 节点中所有使用到的前置节点变量内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) Optional 节点过程数据"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional 输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 执行状态 ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional 耗时(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) optional 总使用 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) optional 总费用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) optional 货币，如 ",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," workflow 执行结束，成功失败同一事件中不同状态",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"workflow_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string)  执行状态 ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional 输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional 耗时(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional 总使用 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 总步数（冗余），默认 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 结束时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS 音频流事件，即：语音合成输出。内容是Mp3格式的音频块，使用 base64 编码后的字符串，播放的时候直接解码即可。(开启自动播放才有此消息)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 语音合成之后的音频块使用 Base64 编码之后的文本内容，播放的时候直接 base64 解码送入播放器即可"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS 音频流结束事件，收到这个事件表示音频流返回结束。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 结束事件是没有音频的，所以这里是空字符串"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 每 10s 一次的 ping 事件，保持连接存活。"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"invalid_param"}),"，传入参数异常"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"app_unavailable"}),"，App 配置不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_not_initialize"}),"，无可用模型凭据配置"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_quota_exceeded"}),"，模型调用额度不足"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"model_currently_not_support"}),"，当前模型不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_request_error"}),"，workflow 执行失败"]}),`
`,e.jsx(n.li,{children:"500，服务内部异常"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/workflows/run",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/run' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Example: file array as an input variable",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "inputs": {
    "{variable_name}":
    [
      {
      "transfer_method": "local_file",
      "upload_file_id": "{upload_file_id}",
      "type": "{document_type}"
      }
    ]
  }
}
`})})}),e.jsx(n.h3,{children:"Blocking Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "workflow_run_id": "djflajgkldjgd",
    "task_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "data": {
        "id": "fdlsjfjejkghjda",
        "workflow_id": "fldjaslkfjlsda",
        "status": "succeeded",
        "outputs": {
          "text": "Nice to meet you."
        },
        "error": null,
        "elapsed_time": 0.875,
        "total_tokens": 3562,
        "total_steps": 8,
        "created_at": 1705407629,
        "finished_at": 1727807631
    }
}
`})})}),e.jsx(n.h3,{children:"Streaming Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})}),e.jsx(s,{title:"File upload sample code",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`import requests
import json

def upload_file(file_path, user):
    upload_url = "https://api.gofy.ai/v1/files/upload"
    headers = {
        "Authorization": "Bearer app-xxxxxxxx",
    }

    try:
        print("上传文件中...")
        with open(file_path, 'rb') as file:
            files = {
                'file': (file_path, file, 'text/plain')  # 确保文件以适当的MIME类型上传
            }
            data = {
                "user": user,
                "type": "TXT"  # 设置文件类型为TXT
            }

            response = requests.post(upload_url, headers=headers, files=files, data=data)
            if response.status_code == 201:  # 201 表示创建成功
                print("文件上传成功")
                return response.json().get("id")  # 获取上传的文件 ID
            else:
                print(f"文件上传失败，状态码: {response.status_code}")
                return None
    except Exception as e:
        print(f"发生错误: {str(e)}")
        return None

def run_workflow(file_id, user, response_mode="blocking"):
    workflow_url = "https://api.gofy.ai/v1/workflows/run"
    headers = {
        "Authorization": "Bearer app-xxxxxxxxx",
        "Content-Type": "application/json"
    }

    data = {
        "inputs": {
            "orig_mail": [{
                "transfer_method": "local_file",
                "upload_file_id": file_id,
                "type": "document"
            }]
        },
        "response_mode": response_mode,
        "user": user
    }

    try:
        print("运行工作流...")
        response = requests.post(workflow_url, headers=headers, json=data)
        if response.status_code == 200:
            print("工作流执行成功")
            return response.json()
        else:
            print(f"工作流执行失败，状态码: {response.status_code}")
            return {"status": "error", "message": f"Failed to execute workflow, status code: {response.status_code}"}
    except Exception as e:
        print(f"发生错误: {str(e)}")
        return {"status": "error", "message": str(e)}

# 使用示例
file_path = "{your_file_path}"
user = "gofyuser"

# 上传文件
file_id = upload_file(file_path, user)
if file_id:
    # 文件上传成功，继续运行工作流
    result = run_workflow(file_id, user)
    print(result)
else:
    print("文件上传失败，无法执行工作流")
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/:workflow_id/run",method:"POST",title:"执行指定版本 Workflow",name:"#Execute-Specific-Workflow"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"执行指定版本的工作流，通过路径参数指定工作流ID。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) Required 工作流ID，用于指定特定版本的工作流"]}),`
`,e.jsx(n.p,{children:"获取方式：在版本历史界面，点击每个版本条目右侧的复制图标即可复制完整的工作流 ID。"}),`
`]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"inputs"}),` (object) Required
允许传入 App 定义的各变量值。
inputs 参数包含了多组键值对（Key/Value pairs），每组的键对应一个特定变量，每组的值则是该变量的具体值。变量可以是文件列表类型。
文件列表类型变量适用于传入文件结合文本理解并回答问题，仅当模型支持该类型文件解析能力时可用。如果该变量是文件列表类型，该变量对应的值应是列表格式，其中每个元素应包含以下内容：`]}),`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 支持类型：",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"})," 具体类型包含：'TXT', 'MD', 'MARKDOWN', 'MDX', 'PDF', 'HTML', 'XLSX', 'XLS', 'VTT', 'PROPERTIES', 'DOC', 'DOCX', 'CSV', 'EML', 'MSG', 'PPTX', 'PPT', 'XML', 'EPUB'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"})," 具体类型包含：'JPG', 'JPEG', 'PNG', 'GIF', 'WEBP', 'SVG'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," 具体类型包含：'MP3', 'M4A', 'WAV', 'WEBM', 'MPGA'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"})," 具体类型包含：'MP4', 'MOV', 'MPEG', 'WEBM'"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," 具体类型包含：其他文件类型"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_method"})," (string)  传递方式:",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"remote_url"}),": 文件地址。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"local_file"}),": 上传文件。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"url"})," 文件地址。（仅当传递方式为 ",e.jsx(n.code,{children:"remote_url"})," 时）。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"upload_file_id"})," 上传文件 ID。（仅当传递方式为 ",e.jsx(n.code,{children:"local_file "}),"时）。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"response_mode"}),` (string) Required
返回响应模式，支持：`]}),`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"streaming"})," 流式模式（推荐）。基于 SSE（",e.jsx(n.strong,{children:e.jsx(n.a,{href:"https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events",children:"Server-Sent Events"})}),"）实现类似打字机输出方式的流式返回。"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"blocking"}),` 阻塞模式，等待执行完毕后返回结果。（请求若流程较长可能会被中断）。
`,e.jsx("i",{children:"由于 Cloudflare 限制，请求会在 100 秒超时无返回后中断。"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
用户标识，用于定义终端用户的身份，方便检索、统计。
由开发者定义规则，需保证用户标识在应用内唯一。API 无法访问 WebApp 创建的会话。`]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"files"})," (array[object]) 可选"]}),`
`]}),`
`,e.jsxs(n.li,{children:[`
`,e.jsxs(n.p,{children:[e.jsx(n.code,{children:"trace_id"}),` (string) Optional
链路追踪ID。适用于与业务系统已有的trace组件打通，实现端到端分布式追踪等场景。如果未指定，系统将自动生成 `,e.jsx(n.code,{children:"trace_id"}),"。支持以下三种方式传递，具体优先级依次为："]}),`
`,e.jsxs(n.ol,{children:[`
`,e.jsxs(n.li,{children:["Header：推荐通过 HTTP Header ",e.jsx(n.code,{children:"X-Trace-Id"})," 传递，优先级最高。"]}),`
`,e.jsxs(n.li,{children:["Query 参数：通过 URL 查询参数 ",e.jsx(n.code,{children:"trace_id"})," 传递。"]}),`
`,e.jsxs(n.li,{children:["Request Body：通过请求体字段 ",e.jsx(n.code,{children:"trace_id"})," 传递（即本字段）。"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.p,{children:["当 ",e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"blocking"}),` 时，返回 CompletionResponse object。
当 `,e.jsx(n.code,{children:"response_mode"})," 为 ",e.jsx(n.code,{children:"streaming"}),"时，返回 ChunkCompletionResponse object 流式序列。"]}),e.jsx(n.h3,{children:"CompletionResponse"}),e.jsxs(n.p,{children:["返回完整的 App 结果，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"application/json"})," 。"]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 执行状态, ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional 输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional 耗时(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional 总使用 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 总步数（冗余），默认 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 结束时间"]}),`
`]}),`
`]}),`
`]}),e.jsx(n.h3,{children:"ChunkCompletionResponse"}),e.jsxs(n.p,{children:["返回 App 输出的流式块，",e.jsx(n.code,{children:"Content-Type"})," 为 ",e.jsx(n.code,{children:"text/event-stream"}),`。
每个流式块均为 data: 开头，块之间以 `,e.jsx(n.code,{children:"\\n\\n"})," 即两个换行符分隔，如下所示："]}),e.jsx(s,{children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`data: {"event": "text_chunk", "workflow_run_id": "b85e5fc5-751b-454d-b14e-dc5f240b0a31", "task_id": "bd029338-b068-4d34-a331-fc85478922c2", "data": {"text": "\\u4e3a\\u4e86", "from_variable_selector": ["1745912968134", "text"]}}\\n\\n
`})})}),e.jsxs(n.p,{children:["流式块中根据 ",e.jsx(n.code,{children:"event"})," 不同，结构也不同，包含以下类型："]}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_started"})," workflow 开始执行",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"workflow_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_started"})," node 开始执行",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"node_started"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) 节点 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_type"})," (string) 节点类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) 节点名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 执行序号，用于展示 Tracing Node 顺序"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) 前置节点 ID，用于画布展示执行路径"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 节点中所有使用到的前置节点变量内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: text_chunk"})," 文本片段",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"text_chunk"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text"})," (string) 文本内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"from_variable_selector"})," (array) 文本来源路径，帮助开发者了解文本是由哪个节点的哪个变量生成的"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: node_finished"})," node 执行结束，成功失败同一事件中不同状态",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"node_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) node 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"node_id"})," (string) 节点 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"index"})," (int) 执行序号，用于展示 Tracing Node 顺序"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"predecessor_node_id"})," (string) optional 前置节点 ID，用于画布展示执行路径"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (object) 节点中所有使用到的前置节点变量内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"process_data"})," (json) Optional 节点过程数据"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional 输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 执行状态 ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional 耗时(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"execution_metadata"})," (json) 元数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) optional 总使用 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_price"})," (decimal) optional 总费用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"currency"})," (string) optional 货币，如 ",e.jsx(n.code,{children:"USD"})," / ",e.jsx(n.code,{children:"RMB"})]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: workflow_finished"})," workflow 执行结束，成功失败同一事件中不同状态",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event"})," (string) 固定为 ",e.jsx(n.code,{children:"workflow_finished"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (object) 详细内容",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string)  执行状态 ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) Optional 输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) Optional 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) Optional 耗时(s)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) Optional 总使用 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 总步数（冗余），默认 0"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 结束时间"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message"})," TTS 音频流事件，即：语音合成输出。内容是Mp3格式的音频块，使用 base64 编码后的字符串，播放的时候直接解码即可。(开启自动播放才有此消息)",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 语音合成之后的音频块使用 Base64 编码之后的文本内容，播放的时候直接 base64 解码送入播放器即可"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: tts_message_end"})," TTS 音频流结束事件，收到这个事件表示音频流返回结束。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，用于请求跟踪和下方的停止响应接口"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"message_id"})," (string) 消息唯一 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"})," (string) 结束事件是没有音频的，所以这里是空字符串"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (int) 创建时间戳，如：1705395332"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"event: ping"})," 每 10s 一次的 ping 事件，保持连接存活。"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"invalid_param"}),"，传入参数异常"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"app_unavailable"}),"，App 配置不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_not_initialize"}),"，无可用模型凭据配置"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"provider_quota_exceeded"}),"，模型调用额度不足"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"model_currently_not_support"}),"，当前模型不可用"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_not_found"}),"，指定的工作流版本未找到"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"draft_workflow_error"}),"，无法使用草稿工作流版本"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_id_format_error"}),"，工作流ID格式错误，需要UUID格式"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"workflow_request_error"}),"，workflow 执行失败"]}),`
`,e.jsx(n.li,{children:"500，服务内部异常"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/workflows/:workflow_id/run",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/{workflow_id}/run' \\
--header 'Authorization: Bearer {api_key}' \\
--header 'Content-Type: application/json' \\
--data-raw '{
  "inputs": ${JSON.stringify(i.inputs)},
  "response_mode": "streaming",
  "user": "abc-123"
}'`}),e.jsx(s,{title:"Example: file array as an input variable",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "inputs": {
    "{variable_name}":
    [
      {
      "transfer_method": "local_file",
      "upload_file_id": "{upload_file_id}",
      "type": "{document_type}"
      }
    ]
  }
}
`})})}),e.jsx(n.h3,{children:"Blocking Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "workflow_run_id": "djflajgkldjgd",
    "task_id": "9da23599-e713-473b-982c-4328d4f5c78a",
    "data": {
        "id": "fdlsjfjejkghjda",
        "workflow_id": "fldjaslkfjlsda",
        "status": "succeeded",
        "outputs": {
          "text": "Nice to meet you."
        },
        "error": null,
        "elapsed_time": 0.875,
        "total_tokens": 3562,
        "total_steps": 8,
        "created_at": 1705407629,
        "finished_at": 1727807631
    }
}
`})})}),e.jsx(n.h3,{children:"Streaming Mode"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-streaming",children:`  data: {"event": "workflow_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "created_at": 1679586595}}
  data: {"event": "node_started", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "created_at": 1679586595}}
  data: {"event": "node_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "node_id": "dfjasklfjdslag", "node_type": "start", "title": "Start", "index": 0, "predecessor_node_id": "fdljewklfklgejlglsd", "inputs": {}, "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "execution_metadata": {"total_tokens": 63127864, "total_price": 2.378, "currency": "USD"},  "created_at": 1679586595}}
  data: {"event": "workflow_finished", "task_id": "5ad4cb98-f0c7-4085-b384-88c403be6290", "workflow_run_id": "5ad498-f0c7-4085-b384-88cbe6290", "data": {"id": "5ad498-f0c7-4085-b384-88cbe6290", "workflow_id": "dfjasklfjdslag", "outputs": {}, "status": "succeeded", "elapsed_time": 0.324, "total_tokens": 63127864, "total_steps": "1", "created_at": 1679586595, "finished_at": 1679976595}}
  data: {"event": "tts_message", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": "qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"}
  data: {"event": "tts_message_end", "conversation_id": "23dd85f3-1a41-4ea0-b7a9-062734ccfaf9", "message_id": "a8bdc41c-13b2-4c18-bfd9-054b9803038c", "created_at": 1721205487, "task_id": "3bf8a0bb-e73b-4690-9e66-4e429bad8ee7", "audio": ""}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/run/:workflow_run_id",method:"GET",title:"获取 Workflow 执行情况",name:"#get-workflow-run-detail"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"根据 workflow 执行 ID 获取 workflow 任务当前执行结果"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run_id"})," (string) workflow 执行 ID，可在流式返回 Chunk 中获取"]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) workflow 执行 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_id"})," (string) 关联的 Workflow ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 执行状态 ",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"inputs"})," (json) 任务输入内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"outputs"})," (json) 任务输出内容"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) 错误原因"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 任务执行总步数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) 任务执行总 tokens"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 任务开始时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 任务结束时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) 耗时 (s)"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"GET",label:"/workflows/run/:workflow_run_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/workflows/run/:workflow_run_id' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "id": "b1ad3277-089e-42c6-9dff-6820d94fbc76",
    "workflow_id": "19eff89f-ec03-4f75-b0fc-897e7effea02",
    "status": "succeeded",
    "inputs": "{\\"sys.files\\": [], \\"sys.user_id\\": \\"abc-123\\"}",
    "outputs": null,
    "error": null,
    "total_steps": 3,
    "total_tokens": 0,
    "created_at": 1705407629,
    "finished_at": 1727807631,
    "elapsed_time": 30.098514399956912
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/tasks/:task_id/stop",method:"POST",title:"停止响应",name:"#stop-generatebacks"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"仅支持流式模式。"}),e.jsx(n.h3,{children:"Path"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"task_id"})," (string) 任务 ID，可在流式返回 Chunk 中获取"]}),`
`]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user"}),` (string) Required
用户标识，用于定义终端用户的身份，必须和发送消息接口传入 user 保持一致。API 无法访问 WebApp 创建的会话。`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"result"}),' (string) 固定返回 "success"']}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(n.h3,{children:"Request Example"}),e.jsx(s,{title:"Request",tag:"POST",label:"/workflows/tasks/:task_id/stop",targetCode:`curl -X POST '${i.appDetail.api_base_url}/workflows/tasks/:task_id/stop' \\
-H 'Authorization: Bearer {api_key}' \\
-H 'Content-Type: application/json' \\
--data-raw '{"user": "abc-123"}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "result": "success"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/files/upload",method:"POST",title:"上传文件",name:"#files-upload"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsxs(n.p,{children:[`上传文件并在发送消息时使用，可实现图文多模态理解。
支持您的工作流程所支持的任何格式。
`,e.jsx("i",{children:"上传的文件仅供当前终端用户使用。"})]}),e.jsx(n.h3,{children:"Request Body"}),e.jsxs(n.p,{children:["该接口需使用  ",e.jsx(n.code,{children:"multipart/form-data"})," 进行请求。"]}),e.jsxs(t,{children:[e.jsx(d,{name:"file",type:"file",children:e.jsx(n.p,{children:"要上传的文件。"})},"file"),e.jsx(d,{name:"user",type:"string",children:e.jsx(n.p,{children:"用户标识，用于定义终端用户的身份，必须和发送消息接口传入 user 保持一致。服务 API 不会共享 WebApp 创建的对话。"})},"user")]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"成功上传后，服务器会返回文件的 ID 和相关信息。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 文件名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"size"})," (int) 文件大小（byte）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"extension"})," (string) 文件后缀"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mime_type"})," (string) 文件 mime-type"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by"})," (uuid) 上传人 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 上传时间"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"no_file_uploaded"}),"，必须提供文件"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"too_many_files"}),"，目前只接受一个文件"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"unsupported_preview"}),"，该文件不支持预览"]}),`
`,e.jsxs(n.li,{children:["400，",e.jsx(n.code,{children:"unsupported_estimate"}),"，该文件不支持估算"]}),`
`,e.jsxs(n.li,{children:["413，",e.jsx(n.code,{children:"file_too_large"}),"，文件太大"]}),`
`,e.jsxs(n.li,{children:["415，",e.jsx(n.code,{children:"unsupported_file_type"}),"，不支持的扩展名，当前只接受文档类文件"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_connection_failed"}),"，无法连接到 S3 服务"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_permission_denied"}),"，无权限上传文件到 S3"]}),`
`,e.jsxs(n.li,{children:["503，",e.jsx(n.code,{children:"s3_file_too_large"}),"，文件超出 S3 大小限制"]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"POST",label:"/files/upload",targetCode:`curl -X POST '${i.appDetail.api_base_url}/files/upload' \\
--header 'Authorization: Bearer {api_key}' \\
--form 'file=@localfile;type=image/[png|jpeg|jpg|webp|gif]' \\
--form 'user=abc-123'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "72fa9618-8f89-4a37-9b33-7e1178a24a67",
  "name": "example.png",
  "size": 1024,
  "extension": "png",
  "mime_type": "image/png",
  "created_by": 123,
  "created_at": 1577836800,
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/end-users/:end_user_id",method:"GET",title:"获取终端用户",name:"#end-user"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"通过终端用户 ID 获取终端用户信息。"}),e.jsxs(n.p,{children:["当其他 API 返回终端用户 ID（例如：上传文件接口返回的 ",e.jsx(n.code,{children:"created_by"}),"）时，可使用该接口查询对应的终端用户信息。"]}),e.jsx(n.h3,{children:"路径参数"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"end_user_id"}),` (uuid) 必需
终端用户 ID。`]}),`
`]}),e.jsx(n.h3,{children:"Response"}),e.jsx(n.p,{children:"返回 EndUser 对象。"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (uuid) ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tenant_id"})," (uuid) 工作空间（Tenant）ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"app_id"})," (uuid) 应用 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 终端用户类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"external_user_id"})," (string) 外部用户 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (boolean) 是否匿名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) 会话 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (string) ISO 8601 时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"updated_at"})," (string) ISO 8601 时间"]}),`
`]}),e.jsx(n.h3,{children:"Errors"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:["404，",e.jsx(n.code,{children:"end_user_not_found"}),"，终端用户不存在"]}),`
`,e.jsx(n.li,{children:"500，内部服务器错误"}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/end-users/:end_user_id",targetCode:`curl -X GET '${i.appDetail.api_base_url}/end-users/6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13' \\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "id": "6ad1ab0a-73ff-4ac1-b9e4-cdb312f71f13",
  "tenant_id": "8c0f3f3a-66b0-4b55-a0bf-8b8e0d6aee7d",
  "app_id": "6c8c3f41-2c6f-4e1b-8f4f-7f11c8f2ad2a",
  "type": "service_api",
  "external_user_id": "abc-123",
  "name": "Alice",
  "is_anonymous": false,
  "session_id": "abc-123",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/workflows/logs",method:"GET",title:"获取 Workflow 日志",name:"#Get-Workflow-Logs"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"倒序返回 workflow 日志"}),e.jsx(n.h3,{children:"Query"}),e.jsxs(t,{children:[e.jsx(d,{name:"keyword",type:"string",children:e.jsx(n.p,{children:"关键字"})},"keyword"),e.jsx(d,{name:"status",type:"string",children:e.jsx(n.p,{children:"执行状态 succeeded/failed/stopped"})},"status"),e.jsx(d,{name:"page",type:"int",children:e.jsx(n.p,{children:"当前页码, 默认1."})},"page"),e.jsx(d,{name:"limit",type:"int",children:e.jsx(n.p,{children:"每页条数, 默认20."})},"limit"),e.jsx(d,{name:"created_by_end_user_session_id",type:"str",children:e.jsxs(n.p,{children:["由哪个endUser创建，例如，",e.jsx(n.code,{children:"abc-123"}),"."]})},"created_by_end_user_session_id"),e.jsx(d,{name:"created_by_account",type:"str",children:e.jsx(n.p,{children:"由哪个邮箱账户创建，例如，lizb@test.com."})},"created_by_account")]}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"page"})," (int) 当前页码"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"limit"})," (int) 每页条数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total"})," (int) 总条数"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"has_more"})," (bool) 是否还有更多数据"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"data"})," (array[object]) 当前页码的数据",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 标识"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"workflow_run"})," (object) Workflow 执行日志",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 标识"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"version"})," (string) 版本"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"status"})," (string) 执行状态，",e.jsx(n.code,{children:"running"})," / ",e.jsx(n.code,{children:"succeeded"})," / ",e.jsx(n.code,{children:"failed"})," / ",e.jsx(n.code,{children:"stopped"})]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"error"})," (string) (可选) 错误"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"elapsed_time"})," (float) 耗时，单位秒"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_tokens"})," (int) 消耗的 token 数量"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"total_steps"})," (int) 执行步骤长度"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 开始时间"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"finished_at"})," (timestamp) 结束时间"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_from"})," (string) 来源"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_role"})," (string) 角色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_account"})," (string) (可选) 帐号"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_by_end_user"})," (object) 用户",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"id"})," (string) 标识"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"type"})," (string) 类型"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"is_anonymous"})," (bool) 是否匿名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"session_id"})," (string) 会话标识"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"created_at"})," (timestamp) 创建时间"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/workflows/logs",targetCode:`curl -X GET '${i.appDetail.api_base_url}/workflows/logs'\\
--header 'Authorization: Bearer {api_key}'`}),e.jsx(n.h3,{children:"Response Example"}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
    "page": 1,
    "limit": 1,
    "total": 7,
    "has_more": true,
    "data": [
        {
            "id": "e41b93f1-7ca2-40fd-b3a8-999aeb499cc0",
            "workflow_run": {
                "id": "c0640fc8-03ef-4481-a96c-8a13b732a36e",
                "version": "2024-08-01 12:17:09.771832",
                "status": "succeeded",
                "error": null,
                "elapsed_time": 1.3588523610014818,
                "total_tokens": 0,
                "total_steps": 3,
                "created_at": 1726139643,
                "finished_at": 1726139644
            },
            "created_from": "service-api",
            "created_by_role": "end_user",
            "created_by_account": null,
            "created_by_end_user": {
                "id": "7f7d9117-dd9d-441d-8970-87e5e7e687a3",
                "type": "service_api",
                "is_anonymous": false,
                "session_id": "abc-123"
            },
            "created_at": 1726139644
        }
    ]
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/info",method:"GET",title:"获取应用基本信息",name:"#info"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取应用的基本信息"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"name"})," (string) 应用名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 应用描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"tags"})," (array[string]) 应用标签"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"mode"})," (string) 应用模式"]}),`
`,e.jsx(n.li,{children:"'author_name' (string) 作者名称"}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/info",targetCode:`curl -X GET '${i.appDetail.api_base_url}/info' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "name": "My App",
  "description": "This is my app.",
  "tags": [
    "tag1",
    "tag2"
  ],
  "mode": "workflow",
  "author_name": "Gofy"
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/parameters",method:"GET",title:"获取应用参数",name:"#parameters"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于进入页面一开始，获取功能开关、输入参数名称、类型及默认值等使用。"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"user_input_form"})," (array[object]) 用户输入表单配置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"text-input"})," (object) 文本输入控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"paragraph"})," (object) 段落文本输入控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"select"})," (object) 下拉控件",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"label"})," (string) 控件展示标签名"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"variable"})," (string) 控件 ID"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"required"})," (bool) 是否必填"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default"})," (string) 默认值"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"options"})," (array[string]) 选项值"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_upload"})," (object) 文件上传配置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"document"}),` (object) 文档设置
当前仅支持文档类型：`,e.jsx(n.code,{children:"txt"}),", ",e.jsx(n.code,{children:"md"}),", ",e.jsx(n.code,{children:"markdown"}),", ",e.jsx(n.code,{children:"pdf"}),", ",e.jsx(n.code,{children:"html"}),", ",e.jsx(n.code,{children:"xlsx"}),", ",e.jsx(n.code,{children:"xls"}),", ",e.jsx(n.code,{children:"docx"}),", ",e.jsx(n.code,{children:"csv"}),", ",e.jsx(n.code,{children:"eml"}),", ",e.jsx(n.code,{children:"msg"}),", ",e.jsx(n.code,{children:"pptx"}),", ",e.jsx(n.code,{children:"ppt"}),", ",e.jsx(n.code,{children:"xml"}),", ",e.jsx(n.code,{children:"epub"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 文档数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image"}),` (object) 图片设置
当前仅支持图片类型：`,e.jsx(n.code,{children:"png"}),", ",e.jsx(n.code,{children:"jpg"}),", ",e.jsx(n.code,{children:"jpeg"}),", ",e.jsx(n.code,{children:"webp"}),", ",e.jsx(n.code,{children:"gif"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 图片数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio"}),` (object) 音频设置
当前仅支持音频类型：`,e.jsx(n.code,{children:"mp3"}),", ",e.jsx(n.code,{children:"m4a"}),", ",e.jsx(n.code,{children:"wav"}),", ",e.jsx(n.code,{children:"webm"}),", ",e.jsx(n.code,{children:"amr"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 音频数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video"}),` (object) 视频设置
当前仅支持视频类型：`,e.jsx(n.code,{children:"mp4"}),", ",e.jsx(n.code,{children:"mov"}),", ",e.jsx(n.code,{children:"mpeg"}),", ",e.jsx(n.code,{children:"mpga"}),"。",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 视频数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom"})," (object) 自定义设置",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"enabled"})," (bool) 是否启用"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"number_limits"})," (int) 自定义数量限制，默认为 3"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"transfer_methods"})," (array[string]) 传输方式列表：",e.jsx(n.code,{children:"remote_url"}),", ",e.jsx(n.code,{children:"local_file"}),"，必须选择一个。"]}),`
`]}),`
`]}),`
`]}),`
`]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"system_parameters"})," (object) 系统参数",`
`,e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"file_size_limit"})," (int) 文档上传大小限制 (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"image_file_size_limit"})," (int) 图片文件上传大小限制（MB）"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"audio_file_size_limit"})," (int) 音频文件上传大小限制 (MB)"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"video_file_size_limit"})," (int) 视频文件上传大小限制 (MB)"]}),`
`]}),`
`]}),`
`]})]}),e.jsxs(r,{sticky:!0,children:[e.jsx(s,{title:"Request",tag:"GET",label:"/parameters",targetCode:` curl -X GET '${i.appDetail.api_base_url}/parameters'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "user_input_form": [
      {
          "paragraph": {
              "label": "Query",
              "variable": "query",
              "required": true,
              "default": ""
          }
      }
  ],
  "file_upload": {
      "image": {
          "enabled": false,
          "number_limits": 3,
          "detail": "high",
          "transfer_methods": [
              "remote_url",
              "local_file"
          ]
      }
  },
  "system_parameters": {
      "file_size_limit": 15,
      "image_file_size_limit": 10,
      "audio_file_size_limit": 50,
      "video_file_size_limit": 100
  }
}
`})})})]})]}),`
`,e.jsx(n.hr,{}),`
`,e.jsx(l,{url:"/site",method:"GET",title:"获取应用 WebApp 设置",name:"#site"}),`
`,e.jsxs(c,{children:[e.jsxs(r,{children:[e.jsx(n.p,{children:"用于获取应用的 WebApp 设置"}),e.jsx(n.h3,{children:"Response"}),e.jsxs(n.ul,{children:[`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"title"})," (string) WebApp 名称"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_type"})," (string) 图标类型，",e.jsx(n.code,{children:"emoji"}),"-表情，",e.jsx(n.code,{children:"image"}),"-图片"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon"})," (string) 图标，如果是 ",e.jsx(n.code,{children:"emoji"})," 类型，则是 emoji 表情符号，如果是 ",e.jsx(n.code,{children:"image"})," 类型，则是图片 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_background"})," (string) hex 格式的背景色"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"icon_url"})," (string) 图标 URL"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"description"})," (string) 描述"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"copyright"})," (string) 版权信息"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"privacy_policy"})," (string) 隐私政策链接"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"custom_disclaimer"})," (string) 自定义免责声明"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"default_language"})," (string) 默认语言"]}),`
`,e.jsxs(n.li,{children:[e.jsx(n.code,{children:"show_workflow_steps"})," (bool) 是否显示工作流详情"]}),`
`]})]}),e.jsxs(r,{children:[e.jsx(s,{title:"Request",tag:"GET",label:"/site",targetCode:`curl -X GET '${i.appDetail.api_base_url}/site' \\
-H 'Authorization: Bearer {api_key}'`}),e.jsx(s,{title:"Response",children:e.jsx(n.pre,{children:e.jsx(n.code,{className:"language-json",children:`{
  "title": "My App",
  "icon_type": "emoji",
  "icon": "😄",
  "icon_background": "#FFEAD5",
  "icon_url": null,
  "description": "This is my app.",
  "copyright": "all rights reserved",
  "privacy_policy": "",
  "custom_disclaimer": "All generated by AI",
  "default_language": "en-US",
  "show_workflow_steps": false,
}
`})})})]})]}),`
`,e.jsx(n.hr,{})]})}function Yn(i={}){const{wrapper:n}=i.components||{};return n?e.jsx(n,{...i,children:e.jsx(we,{...i})}):we(i)}const Zn=({toc:i,activeSection:n,isTocExpanded:a,onToggle:o,onItemClick:h})=>{const{t:x}=Te();return a?e.jsxs("nav",{className:"toc flex max-h-[calc(100vh-150px)] w-full flex-col overflow-hidden rounded-xl border-[0.5px] border-components-panel-border bg-background-default-hover shadow-xl",children:[e.jsxs("div",{className:"relative z-10 flex items-center justify-between border-b border-components-panel-border-subtle bg-background-default-hover px-4 py-2.5",children:[e.jsx("span",{className:"text-xs font-medium uppercase tracking-wide text-text-tertiary",children:x("develop.toc",{ns:"appApi"})}),e.jsx("button",{type:"button",onClick:()=>o(!1),className:"group flex h-6 w-6 items-center justify-center rounded-md transition-colors hover:bg-state-base-hover","aria-label":"Close",children:e.jsx("span",{className:"i-ri-close-line h-3 w-3 text-text-quaternary transition-colors group-hover:text-text-secondary"})})]}),e.jsx("div",{className:"from-components-panel-border-subtle/20 pointer-events-none absolute left-0 right-0 top-[41px] z-10 h-2 bg-gradient-to-b to-transparent"}),e.jsx("div",{className:"pointer-events-none absolute left-0 right-0 top-[43px] z-10 h-3 bg-gradient-to-b from-background-default-hover to-transparent"}),e.jsx("div",{className:"relative flex-1 overflow-y-auto px-3 py-3 pt-1",children:i.length===0?e.jsx("div",{className:"px-2 py-8 text-center text-xs text-text-quaternary",children:x("develop.noContent",{ns:"appApi"})}):e.jsx("ul",{className:"space-y-0.5",children:i.map(j=>{const _=n===j.href.replace("#","");return e.jsx("li",{children:e.jsxs("a",{href:j.href,onClick:y=>h(y,j),className:R("group relative flex items-center rounded-md px-3 py-2 text-[13px] transition-all duration-200",_?"bg-state-base-hover font-medium text-text-primary":"text-text-tertiary hover:bg-state-base-hover hover:text-text-secondary"),children:[e.jsx("span",{className:R("mr-2 h-1.5 w-1.5 rounded-full transition-all duration-200",_?"scale-100 bg-text-accent":"scale-75 bg-components-panel-border")}),e.jsx("span",{className:"flex-1 truncate",children:j.text})]})},j.href)})})}),e.jsx("div",{className:"pointer-events-none absolute bottom-0 left-0 right-0 z-10 h-4 rounded-b-xl bg-gradient-to-t from-background-default-hover to-transparent"})]}):e.jsx("button",{type:"button",onClick:()=>o(!0),className:"group flex h-11 w-11 items-center justify-center rounded-full border-[0.5px] border-components-panel-border bg-components-panel-bg shadow-lg transition-all duration-150 hover:bg-background-default-hover hover:shadow-xl","aria-label":"Open table of contents",children:e.jsx("span",{className:"i-ri-list-unordered h-5 w-5 text-text-tertiary transition-colors group-hover:text-text-secondary"})})},es={[$.CHAT]:{zh:ye,ja:fe,en:ge},[$.AGENT_CHAT]:{zh:ye,ja:fe,en:ge},[$.ADVANCED_CHAT]:{zh:Kn,ja:Vn,en:Fn},[$.WORKFLOW]:{zh:Yn,ja:Jn,en:Qn},[$.COMPLETION]:{zh:Hn,ja:$n,en:Wn}},ns=(i,n)=>{if(!i)return null;const a=es[i];if(!a)return null;const o=$e(n);return a[o]??a.en??null},ss=({appDetail:i})=>{var g,k;const n=Qe(),{theme:a}=Je(),{toc:o,isTocExpanded:h,setIsTocExpanded:x,activeSection:j,handleTocClick:_}=Rn({appDetail:i,locale:n}),y=((k=(g=i==null?void 0:i.model_config)==null?void 0:g.configs)==null?void 0:k.prompt_variables)??[],m=y.reduce((b,f)=>(b[f.key]=f.name||"",b),{}),p=u.useMemo(()=>ns(i==null?void 0:i.mode,n),[i==null?void 0:i.mode,n]);return e.jsxs("div",{className:"flex",children:[e.jsx("div",{className:`fixed right-20 top-32 z-10 transition-all duration-150 ease-out ${h?"w-[280px]":"w-11"}`,children:e.jsx(Zn,{toc:o,activeSection:j,isTocExpanded:h,onToggle:x,onItemClick:_})}),e.jsx("article",{className:R("prose-xl prose",a===We.dark&&"prose-invert"),children:p&&e.jsx(p,{appDetail:i,variables:y,inputs:m})})]})},is=({appId:i})=>{const n=Fe(a=>a.appDetail);return n?e.jsxs("div",{"data-testid":"develop-main",className:"relative flex h-full flex-col overflow-hidden",children:[e.jsxs("div",{className:"flex shrink-0 items-center justify-between border-b border-solid border-b-divider-regular px-6 py-2",children:[e.jsx("div",{className:"text-lg font-medium text-text-primary"}),e.jsx(In,{apiBaseUrl:n.api_base_url,appId:i})]}),e.jsx("div",{className:"grow overflow-auto px-4 py-4 sm:px-10",children:e.jsx(ss,{appDetail:n})})]}):e.jsx("div",{className:"flex h-full items-center justify-center bg-background-default",children:e.jsx(He,{})})},Us=async i=>{const n=await i.params,{appId:a}=n;return e.jsx(is,{appId:a})};export{Us as default};
