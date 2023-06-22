function Decode(script) {
 try {
   /**
    *
    * @param {string} h
    * @param {number} u
    * @param {string} n
    * @param {number} t
    * @param {number} e
    * @param {number} r
    * @returns {string}
    */
   const decodeSnap = (h, u, n, t, e, r) => {
     /**
      *
      * @param {string} d
      * @param {number} e
      * @param {number} f
      * @returns {string}
      */
     function chip(d, e, f) {
       var g =
           "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/".split(
               ""
           );
       var h = g.slice(0, e);
       var i = g.slice(0, f);
       var j = d
           .split("")
           .reverse()
           .reduce(function (a, b, c) {
             if (h.indexOf(b) !== -1)
               return (a += h.indexOf(b) * Math.pow(e, c));
             return a;
           }, 0);
       var k = "";
       while (j > 0) {
         k = i[j % f] + k;
         j = (j - (j % f)) / f;
       }
       return k || "0";
     }

     r = "";
     for (var i = 0, len = h.length; i < len; i++) {
       var s = "";
       while (h[i] !== n[e]) {
         s += h[i];
         i++;
       }
       for (var j = 0; j < n.length; j++) {
         s = s.replace(new RegExp(n[j], "g"), j);
       }
       r += String.fromCharCode(chip(s, e, 10) - t);
     }
     return decodeURIComponent(escape(r));
   };

   const passRegExp = /\}eval\(function/g.exec(script);

   if (!passRegExp) {
     throw new Error("[404] Could not find executable script.");
   }

   const paramRegExp = /escape\(r\)\)\}\((.*?)\)\)/
       .exec(script)[1]
       .split(",")
       .map((v) => (v.includes('"') ? v.slice(1, -1) : parseInt(v)));

   const [h, u, n, t, e, r] = paramRegExp;

   return decodeSnap(h, u, n, t, e, r);
 } catch (error) {
   throw error;
 }
}