package erogs

import (
	"fmt"
	"strings"

	internalerrors "kurohelper/errors"
	"kurohelper/utils"
)

func buildFuzzySearchCreatorSQL(search string) (string, error) {
	search = strings.ReplaceAll(search, "'", "''")
	result := "%"
	if utils.IsAllEnglish(search) {
		result += search + "%"
	} else {
		for _, r := range search {
			result += string(r) + "%"
		}
	}

	if strings.TrimSpace(search) == "" {
		return "", internalerrors.ErrSearchNoContent
	}

	return fmt.Sprintf(`
SELECT row_to_json(c)
FROM (
    SELECT
        cr.id,
        cr.name,
        cr.twitter_username,
        cr.blog,
        cr.pixiv,
        (
            SELECT json_agg(game_data)
            FROM (
                SELECT
                    g.gamename,
                    g.sellday,
                    g.median,
                    g.count_all,
                    (
                        SELECT json_agg(
                            json_build_object(
                                'shubetu', s2.shubetu,
                                'shubetu_detail', s2.shubetu_detail,
                                'shubetu_detail_name', s2.shubetu_detail_name
                            )
                        )
                        FROM shokushu s2
                        WHERE s2.creater = cr.id
                          AND s2.game = g.id
                    ) AS shokushu
                FROM gamelist g
                WHERE EXISTS (
                    SELECT 1 
                    FROM shokushu s3
                    WHERE s3.creater = cr.id
                      AND s3.game = g.id
                )
                GROUP BY g.id, g.gamename, g.sellday, g.median, g.count_all
            ) AS game_data
        ) AS games
    FROM createrlist cr
    WHERE cr.name ILIKE '%s'
    LIMIT 1
) AS c;`, result), nil
}

func buildFuzzySearchMusicSQL(search string) (string, error) {
	search = strings.ReplaceAll(search, "'", "''")
	result := "%"
	if utils.IsAllEnglish(search) {
		result += search + "%"
	} else {
		for _, r := range search {
			result += string(r) + "%"
		}
	}

	if strings.TrimSpace(search) == "" {
		return "", internalerrors.ErrSearchNoContent
	}

	return fmt.Sprintf(`
WITH singer_agg AS (
    SELECT music AS music_id,
           json_agg(
               json_build_object(
                   'singer_name', c.name
               )
           ) AS singers
    FROM singer
    JOIN createrlist AS c ON c.id = singer.creater
    GROUP BY music
),
lyrics_agg AS (
    SELECT music AS music_id,
           json_agg(
               json_build_object(
                   'lyric_name', c.name
               )
           ) AS lyrics
    FROM lyrics
    JOIN createrlist AS c ON c.id = lyrics.creater
    GROUP BY music
),
arrangement_agg AS (
    SELECT music AS music_id,
           json_agg(
               json_build_object(
                   'arrangement_name', c.name
               )
           ) AS arrangements
    FROM arrangement
    JOIN createrlist AS c ON c.id = arrangement.creater
    GROUP BY music
),
composition_agg AS (
    SELECT music AS music_id,
           json_agg(
               json_build_object(
                   'composition_name', c.name
               )
           ) AS compositions
    FROM composition
    JOIN createrlist AS c ON c.id = composition.creater
    GROUP BY music
),
gamelist_agg AS (
    SELECT gm.music AS music_id,
           json_agg(
               json_build_object(
                   'game_name', g.gamename,
                   'category', gm.category
               )
           ) AS games
    FROM game_music gm
    JOIN gamelist g ON g.id = gm.game
    GROUP BY gm.music
),
musicitemlist_agg AS (
    SELECT music AS music_id,
           json_agg(
               json_build_object(
                   'album_name', mi.name
               )
           ) AS albums
    FROM musicitem_music mim
    JOIN musicitemlist mi ON mi.id = mim.musicitem
    GROUP BY music
),
usermusic_tokuten_agg AS (
    SELECT music AS music_id,
           ROUND(AVG(LEAST(tokuten, 100)),2) AS avg_tokuten,
           COUNT(uid) AS tokuten_count
    FROM usermusic_tokuten
    GROUP BY music
)
SELECT json_agg(row_to_json(t))
FROM (
    SELECT m.id AS music_id,
           m.name AS songname,
           m.playtime,
           m.releasedate,
           ut.avg_tokuten,
           ut.tokuten_count,
           s.singers,
           l.lyrics,
           a.arrangements,
           comp.compositions,
           g.games,
           mi.albums
    FROM musiclist m
    LEFT JOIN singer_agg s ON s.music_id = m.id
    LEFT JOIN lyrics_agg l ON l.music_id = m.id
    LEFT JOIN arrangement_agg a ON a.music_id = m.id
    LEFT JOIN composition_agg comp ON comp.music_id = m.id
    LEFT JOIN gamelist_agg g ON g.music_id = m.id
    LEFT JOIN musicitemlist_agg mi ON mi.music_id = m.id
    LEFT JOIN usermusic_tokuten_agg ut ON ut.music_id = m.id 
    WHERE m.name ILIKE '%s'
    LIMIT 2000
) t;
`, result), nil
}
