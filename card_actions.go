//go:build js && wasm

package main

import (
	"strings"
	"syscall/js"
)

var activePaletteRefresh func()

var paletteGroups = map[string]string{
	"顔":    "😀 😃 😄 😁 😆 😅 😂 🤣 😊 😇 🙂 🙃 😉 😌 😍 🥰 😘 😗 😙 😚 😋 😛 😜 🤪 😝 🤑 🤗 🤭 🤫 🤔 🤐 🤨 😐 😑 😶 😏 😒 🙄 😬 😮‍💨 🤥 😌 😔 😪 🤤 😴 😷 🤒 🤕 🤢 🤮 🤧 🥵 🥶 🥴 😵 🤯 🤠 🥳 🥸 😎 🤓 🧐 😕 😟 🙁 ☹️ 😮 😯 😲 😳 🥺 😦 😧 😨 😰 😥 😢 😭 😱 😖 😣 😞 😓 😩 😫 🥱 😤 😡 😠 🤬 😈 👿 💀 ☠️ 💩 🤡 👻 👽 🤖",
	"手・人":  "👋 🤚 🖐️ ✋ 🖖 👌 🤌 🤏 ✌️ 🤞 🫰 🤟 🤘 🤙 👈 👉 👆 🖕 👇 ☝️ 👍 👎 ✊ 👊 🤛 🤜 👏 🙌 👐 🤲 🤝 🙏 ✍️ 💪 🫶 👀 👁️ 👄 🧠 🧑 👶 👧 👦 👩 👨 🧓 👮 🧑‍💻 🧑‍🎨 🧑‍🚀 💃 🕺 🏃 🚶 🧘",
	"自然":   "🐶 🐱 🐭 🐹 🐰 🦊 🐻 🐼 🐨 🐯 🦁 🐮 🐷 🐸 🐵 🐔 🐧 🐦 🐤 🦆 🦉 🦄 🐝 🦋 🐌 🐞 🐢 🐙 🐟 🐬 🐳 🐋 🦀 🦐 🐉 🌵 🌲 🌳 🌴 🌱 🌿 ☘️ 🍀 🍁 🍂 🍄 🌸 🌼 🌻 🌹 🌷 🌺 🌈 ☀️ 🌤️ ⛅ 🌧️ ⛈️ ❄️ ☃️ 🌊 🔥 ⭐ 🌙 🌍",
	"食べ物":  "🍎 🍐 🍊 🍋 🍌 🍉 🍇 🍓 🫐 🍈 🍒 🍑 🥭 🍍 🥥 🥝 🍅 🥑 🍆 🥕 🌽 🥦 🍞 🥐 🥖 🧀 🍖 🍗 🍔 🍟 🍕 🌭 🥪 🌮 🍣 🍱 🍛 🍜 🍝 🍙 🍚 🍘 🍡 🍦 🍩 🍪 🎂 🍰 🍫 🍬 ☕ 🍵 🧃 🍺 🍻 🍷 🥂",
	"記号・物": "❤️ 🧡 💛 💚 💙 💜 🖤 🤍 🤎 💔 ❣️ 💕 💞 💓 💗 💖 💘 💝 💯 💢 💥 💫 💦 💤 ✨ 🎉 🎊 🎁 🎈 🎵 🎶 🎸 🎮 ⚽ 🏀 🏆 🚗 🚀 ✈️ 🏠 💻 📱 📷 💡 📚 ✏️ ✅ ❌ ⭕ ❓ ❗ ⚠️ ♻️ 🌟 🇯🇵 🇺🇸",
}

func buildReactionActions(actions js.Value, ev Event) {
	button := doc.Call("createElement", "button")
	if selectedReaction.Code != "" {
		img := doc.Call("createElement", "img")
		img.Set("src", selectedReaction.URL)
		img.Set("alt", selectedReaction.Content)
		img.Set("className", "default-reaction-image")
		button.Call("appendChild", img)
	} else {
		button.Set("textContent", selectedReaction.Content)
	}
	button.Set("ariaLabel", "デフォルトのリアクションを送る")
	button.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
		publishReaction(ev.ID, ev.Pubkey, selectedReaction.Content, selectedReaction.Code, selectedReaction.URL)
		return nil
	}))
	actions.Call("appendChild", button)
	toggle := doc.Call("createElement", "button")
	toggle.Set("textContent", "+")
	toggle.Set("ariaLabel", "絵文字を選ぶ")
	toggle.Call("setAttribute", "aria-expanded", "false")
	actions.Call("appendChild", toggle)
	choices := doc.Call("createElement", "div")
	choices.Set("className", "emoji-choices")
	choices.Set("hidden", true)
	actions.Call("appendChild", choices)
	var show func(string)
	show = func(tab string) {
		activePaletteRefresh = func() {
			if !choices.Get("hidden").Bool() {
				show(tab)
			}
		}
		choices.Set("textContent", "")
		tabs := doc.Call("createElement", "div")
		tabs.Set("className", "emoji-tabs")
		choices.Call("appendChild", tabs)
		for _, name := range []string{"recent", "顔", "手・人", "自然", "食べ物", "記号・物", "カスタム"} {
			name := name
			b := doc.Call("createElement", "button")
			b.Set("textContent", name)
			b.Call("setAttribute", "aria-pressed", name == tab)
			b.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any { show(name); return nil }))
			tabs.Call("appendChild", b)
		}
		grid := doc.Call("createElement", "div")
		grid.Set("className", "emoji-grid")
		choices.Call("appendChild", grid)
		items := []recentEmoji{}
		if tab == "recent" {
			items = loadRecent()
		} else if tab == "カスタム" {
			seen := map[string]bool{}
			for _, t := range availableEmojiTags() {
				if len(t) >= 3 && t[0] == "emoji" && safeHTTPS(t[2]) && !seen[t[1]+"\x00"+t[2]] {
					seen[t[1]+"\x00"+t[2]] = true
					items = append(items, recentEmoji{":" + t[1] + ":", t[1], t[2]})
				}
			}
		} else {
			for _, v := range strings.Fields(paletteGroups[tab]) {
				items = append(items, recentEmoji{Content: v})
			}
		}
		for _, item := range items {
			item := item
			if item.Code == "" && !validEmoji(item.Content) || item.Code != "" && !safeHTTPS(item.URL) {
				continue
			}
			b := doc.Call("createElement", "button")
			b.Set("title", item.Content)
			b.Set("ariaLabel", item.Content)
			if item.Code != "" {
				img := doc.Call("createElement", "img")
				img.Set("src", item.URL)
				img.Set("alt", item.Content)
				img.Set("loading", "lazy")
				b.Call("appendChild", img)
			} else {
				b.Set("textContent", item.Content)
			}
			b.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
				publishReaction(ev.ID, ev.Pubkey, item.Content, item.Code, item.URL)
				choices.Set("hidden", true)
				toggle.Call("setAttribute", "aria-expanded", "false")
				return nil
			}))
			grid.Call("appendChild", b)
		}
		if len(items) == 0 {
			grid.Set("textContent", "絵文字はまだありません")
		}
	}
	toggle.Call("addEventListener", "click", callback(func(_ js.Value, _ []js.Value) any {
		hidden := !choices.Get("hidden").Bool()
		choices.Set("hidden", hidden)
		toggle.Call("setAttribute", "aria-expanded", !hidden)
		if !hidden {
			show("recent")
		}
		return nil
	}))
}
