// ==UserScript==
// @name         ActorPlus
// @description:zh-cn 隐藏没有头像的演员和制作人员
// @version      c85f7ce
// @author       @newday-life
// @github       https://github.com/newday-life/emby-front-end-mod/blob/main/actorPlus/actorPlus.js
// ==/UserScript==

(function () {
    "use strict";
    /* page item.Type "Person" "Movie" "Series" "Season" "Episode" "BoxSet" "video-osd" so. */
    const show_pages = ["Movie", "Series", "Episode", "Season", "video-osd"];
    var item, paly_mutation;

    /* document.addEventListener("itemshow", function (e) {
        item = e.detail.item;
        // if (showFlag() && item.People) {
        //     item.People = item.People.filter(p => p.PrimaryImageTag);
        // }
    }); */
    document.addEventListener("viewbeforeshow", function (e) {
        paly_mutation?.disconnect();
        if (e.detail.path === "/item" || e.detail.type === "video-osd") {
            if (!e.detail.isRestored) {
                const mutation = new MutationObserver(async function () {
                    item = e.target.controller?.currentItem || e.target.controller?.videoOsd?.currentItem || e.target.controller?.currentPlayer?.streamInfo?.item;
                    if (item) {
                        mutation.disconnect();
                        if (showFlag()) {
                            if (!item.People) {
                                item = await ApiClient.getItem(ApiClient.getCurrentUserId(), item.Id);
                            }
                            if (item.People.length === 0) {
                                return;
                            }
                            e.detail.path === "/item" && (e.target.querySelector(".peopleItemsContainer").fetchData = filterPeople.bind(item));
                            /* item.People = item.People.filter(p => p.PrimaryImageTag);
                            e.detail.type === "video-osd" && setTimeout(() => {
                                e.target.controller.videoOsd && (e.target.controller.videoOsd.currentItem.People = item.People);
                                e.target.controller.currentItem && (e.target.controller.currentItem.People = item.People);
                            }, 1000); */
                            if (e.detail.type === "video-osd") {
                                paly_mutation = new MutationObserver(function () {
                                    let itemsContainer = e.target.querySelector('[data-index="2"].videoosd-itemstab .itemsContainer');
                                    if (itemsContainer) {
                                        paly_mutation.disconnect();
                                        itemsContainer.fetchData = filterPeople.bind(item);
                                    }

                                });
                                paly_mutation.observe(e.target.querySelector('[data-index="2"].videoosd-itemstab'), {
                                    childList: true,
                                    characterData: true,
                                    subtree: true,
                                });
                            }
                        }
                    }
                });
                mutation.observe(document.body, {
                    childList: true,
                    characterData: true,
                    subtree: true,
                });
            } else {
                item = e.target.controller.currentItem || e.target.controller.videoOsd?.currentItem;
            }
        }
    });
    function filterPeople(query) {
        var serverId = item.ServerId,
            people = (item.People || []).filter(function (p) {
                return p.PrimaryImageTag && (p.ServerId = serverId, "Person" !== p.Type && (p.PersonType = p.Type, p.Type = "Person"), !0)
            }),
            totalRecordCount = people.length;
        return query && (people = people.slice(query.StartIndex || 0),
            query.Limit && people.length > query.Limit && (people.length = query.Limit)),
            Promise.resolve({
                Items: people,
                TotalRecordCount: totalRecordCount
            })
    }
    function showFlag() {
        for (let show_page of show_pages) {
            if (item.Type == show_page) {
                return true;
            }
        }
        return false;
    }

})();