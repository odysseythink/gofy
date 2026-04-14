const l=(t,o,n)=>{let e=`mailto:${t}`;return o&&(e+=`?subject=${encodeURIComponent(o)}`),n&&(e+=`&body=${encodeURIComponent(n)}`),e},s=(t,o,n,e)=>{const r=`Technical Support Request ${o} ${t}`,i=`
    Please do not remove the following information:
    -----------------------------------------------
    Current Plan: ${o}
    Account: ${t}
    Version: ${n}
    Platform:
    Problem Description:
  `;return l(e||"support@gofy.ai",r,i)};export{s as m};
