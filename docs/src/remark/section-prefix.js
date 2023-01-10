const visit = require('unist-util-visit');
const diffParser = require('parse-diff');

const plugin = (options) => {
	const transformer = async (ast) => {
		visit(ast, 'code', (node) => {
			if (node.lang == 'patch') {
				diff=diffParser(node.value)
				node.lang = 'go'
				// node.value = node.value.replace('-', '// remove-next-line\n');
				// node.value = node.value.replace('+', '// highlight-next-line\n');
				// console.log(diff)
				diff=diff[0]
				node.meta='title="'+diff.to+'"';
				let s="";
				for (i=0; i<diff.chunks.length;i++){
					// console.log("chunk",diff.chunks[i])
					content=diff.chunks[i].content
					content=content.substring(content.lastIndexOf('@@ ')+3)
					s+=content+"\n"
					if (diff.chunks[i].oldStart!=diff.chunks[i].newStart){
						s+="\n// ...\n\n"
					}
					for (j=0; j<diff.chunks[i].changes.length;j++){
						switch (diff.chunks[i].changes[j].type) {
							case 'add':
								s+="// highlight-next-line\n"
								break
							case 'del':
								s+="// remove-next-line\n"
								break
						}
						content=diff.chunks[i].changes[j].content.substring(1);
						s+=content+"\n"	
					}
					s+="\n// ...\n\n"	
				}
				node.value=s
			}
		});
	};
	return transformer;
};

module.exports = plugin;
