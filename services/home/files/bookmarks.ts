#!/usr/bin/env bun

const linkdingUrl = "https://linkding.abchome.work.gd";
const apiKey = "1a30142453a89bcc7b67b94a4f30f151b50fd476";


// 1. Giant map of domain/tag keys to their target bundled tags
const tagMapping = {
	// Programming & Development
	frontend: ["programming"],
	"github.com": ["programming", "github-repo"],
	"go.dev": ["programming", "backend", "go"],
	"gobyexample.com": ["programming", "backend", "go"],
	"supabase.com": ["programming", "backend", "databases"],
	"pocketbase.io": ["programming", "backend"],
	"sql-practice.com": ["programming", "databases", "learning-to-code"],
	"mystery.knightlab.com": ["programming", "databases", "learning-to-code"],
	"crontab.guru": ["programming", "scripting", "sysadmin"],
	"regexr.com": ["programming", "scripting"],
	"jqlang.org": ["programming", "scripting"],
	"spoj.com": ["programming", "learning-to-code"],
	"codeforces.com": ["programming", "learning-to-code"],

	// Security & System Admin
	"kali.org": ["cybersecurity", "sysadmin"],
	"exploit-db.com": ["cybersecurity"],
	"overthewire.org": ["cybersecurity"],
	"tryhackme.com": ["cybersecurity"],
	"gitleaks.io": ["cybersecurity", "devops"],
	"shodan.io": ["sysadmin", "networking"],
	"firebog.net": ["sysadmin", "networking"],

	// Entertainment & Sports
	sportsurge: ["streaming-replays", "sports", "live-sports"],
	score808: ["streaming-replays", "sports", "live-sports"],
	motogp: ["streaming-replays", "sports", "motorsports"],
	"f1fullraces.com": ["streaming-replays", "sports", "motorsports"],
	"fullmatchsports.cc": ["streaming-replays", "sports"],
	watchwrestling: ["streaming-replays", "sports", "wrestling"],
	animepahe: ["movies-and-tv", "streaming-replays"],
	"themoviedb.org": ["movies-and-tv"],
	"opensubtitles.org": ["movies-and-tv", "subtitles"],

	// Piracy & Repacks
	"fitgirl-repacks.site": ["torrents", "piracy-indexes", "pc-games"],
	"dodi-repacks.site": ["torrents", "piracy-indexes", "pc-games"],
	"1337x.to": ["torrents", "piracy-indexes"],
	"yts.mx": ["torrents", "piracy-indexes", "movies-and-tv"],
	"eztvx.to": ["torrents", "piracy-indexes", "movies-and-tv"],
	"annas-archive.li": ["open-directories", "direct-downloads", "e-books"],
	"z-library": ["open-directories", "direct-downloads", "e-books"],
	"oceanofpdf.com": ["open-directories", "direct-downloads", "e-books"],
	"vadapav.mov": ["open-directories", "direct-downloads"],
	moviesmod: ["open-directories", "direct-downloads", "movies-and-tv"],

	// Tools & Utilities
	ffmpeg: ["media-tools"],
	y2mate: ["media-tools", "youtube-downloaders"],
	"it-tools": ["dev-utilities"],
	scrcpy: ["dev-utilities", "sysadmin"],
	"massgrave.dev": ["dev-utilities", "windows"],
};

async function bundleBookmarks() {
	const headers = {
		Authorization: `Token ${apiKey}`,
		"Content-Type": "application/json",
	};

	try {
		// Fetch existing entries
		console.log("Fetching bookmarks from Linkding...");
		const res = await fetch(
			`${linkdingUrl.replace(/\/$/, "")}/api/bookmarks/?limit=300`,
			{ headers },
		);

		if (!res.ok) throw new Error(`Fetch error: ${res.statusText}`);
		const data = await res.json();
		const bookmarks = data.results || [];

		// Loop through matches and patch payloads
		for (const bookmark of bookmarks) {
			const url = bookmark.url.toLowerCase();
			const currentTags = bookmark.tag_names || [];
			const updatedTagsSet = new Set(currentTags);

			// 1. Match against existing tags
			currentTags.forEach((tag) => {
				if (tagMapping[tag]) {
					tagMapping[tag].forEach((t) => updatedTagsSet.add(t));
				}
			});

			// 2. Match against domain keywords inside URL
			Object.keys(tagMapping).forEach((key) => {
				if (url.includes(key)) {
					tagMapping[key].forEach((t) => updatedTagsSet.add(t));
				}
			});

			// If tags changed, execute PATCH update
			if (updatedTagsSet.size !== currentTags.length) {
				const newTagsList = Array.from(updatedTagsSet);
				const patchUrl = `${linkdingUrl.replace(/\/$/, "")}/api/bookmarks/${bookmark.id}/`;

				// Only provide field being updated
				const payload = { tag_names: newTagsList };

				const patchRes = await fetch(patchUrl, {
					method: "PATCH",
					headers: headers,
					body: JSON.stringify(payload),
				});

				if (patchRes.ok) {
					console.log(
						`✓ Updated: "${bookmark.title}" -> ${JSON.stringify(newTagsList)}`,
					);
				} else {
					console.error(
						`✗ Failed updating ${bookmark.title}: ${patchRes.status}`,
					);
				}
			}
		}
		console.log("Tag bundling synchronization finished.");
	} catch (error) {
		console.error("Pipeline failure:", error);
	}
}

bundleBookmarks();

// // fetch all bookmarks
// async function allbookmarks() {
// 	try {
// 		const res = await fetch(`${linkdingrl}/api/bookmarks/?limit=300`, {
// 			headers: {
// 				Authorization: `Token ${apiKey}`,
// 			},
// 		});
// 		if (!res.ok) {
// 			throw new Error(res.statusText);
// 		}

// 		const bookmarks = await res.json();

// 		console.log(bookmarks);
// 	} catch (error) {
// 		console.error(error);
// 	}
// }

// allbookmarks();
