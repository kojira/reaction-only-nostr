const viewer = document.getElementById('reactionViewer');
const viewerContent = document.getElementById('reactionViewerContent');
let viewerTrigger;
document.getElementById('feed').addEventListener('click', event => {
  const reaction = event.target.closest('button.reaction');
  if (!reaction) return;
  viewerTrigger = reaction;
  viewerContent.replaceChildren(...Array.from(reaction.childNodes, node => node.cloneNode(true)));
  viewer.showModal();
});
document.getElementById('reactionViewerClose').addEventListener('click', () => viewer.close());
viewer.addEventListener('click', event => { if (event.target === viewer) viewer.close(); });
viewer.addEventListener('close', () => { viewerContent.replaceChildren(); if (viewerTrigger?.isConnected) viewerTrigger.focus(); });
