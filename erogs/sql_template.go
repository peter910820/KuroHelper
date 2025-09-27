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
           STRING_AGG(c.name, ',') AS singer_name 
    FROM singer
    JOIN createrlist AS c ON c.id = singer.creater
    GROUP BY music
),
lyrics_agg AS (
    SELECT music AS music_id,
           STRING_AGG(c.name, ',') AS lyric_name 
    FROM lyrics
    JOIN createrlist AS c ON c.id = lyrics.creater
    GROUP BY music
),
arrangement_agg AS (
    SELECT music AS music_id,
           STRING_AGG(c.name, ',') AS arrangement_name
    FROM arrangement
    JOIN createrlist AS c ON c.id = arrangement.creater
    GROUP BY music
),
composition_agg AS (
    SELECT music AS music_id,
           STRING_AGG(c.name, ',') AS composition_name 
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
           ) AS game_categories
    FROM game_music gm
    JOIN gamelist g ON g.id = gm.game
    GROUP BY gm.music
),
musicitemlist_agg AS (
    SELECT music AS music_id,
           STRING_AGG(mi.name, ',') AS album_name 
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
SELECT row_to_json(t)
FROM (
    SELECT m.id AS music_id,
           m.name AS musicname,
           m.playtime,
           m.releasedate,
           ut.avg_tokuten,
           ut.tokuten_count,
           COALESCE(s.singer_name, '無') AS singer_name,
           COALESCE(l.lyric_name, '無') AS lyric_name,
           COALESCE(a.arrangement_name, '無') AS arrangement_name,
           COALESCE(comp.composition_name, '無') AS composition_name,
           g.game_categories,
           COALESCE(mi.album_name, '') AS album_name
    FROM musiclist m
    LEFT JOIN singer_agg s ON s.music_id = m.id
    LEFT JOIN lyrics_agg l ON l.music_id = m.id
    LEFT JOIN arrangement_agg a ON a.music_id = m.id
    LEFT JOIN composition_agg comp ON comp.music_id = m.id
    LEFT JOIN gamelist_agg g ON g.music_id = m.id
    LEFT JOIN musicitemlist_agg mi ON mi.music_id = m.id
    LEFT JOIN usermusic_tokuten_agg ut ON ut.music_id = m.id 
    WHERE m.name ILIKE '%s'
    ORDER BY ut.avg_tokuten DESC NULLS LAST
    LIMIT 1
)t;
`, result), nil
}

func buildFuzzySearchGameSQL(search string) (string, error) {
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
WITH shokushu_agg AS (
    SELECT s.game AS game_id,
           json_agg(
               json_build_object(
                   'shubetu_type', s.shubetu,
                    'creater_name', c.name,
                    'shubetu_detail_type', s.shubetu_detail,
                    'shubetu_detail_name', s.shubetu_detail_name
                       
               )
           ) AS shubetu_detail
    FROM shokushu s
    LEFT JOIN createrlist c ON c.id = s.creater
    WHERE s.shubetu != 7
    GROUP BY s.game
)
SELECT row_to_json(t)
FROM (
    SELECT g.id,
           b.brandname,
           g.gamename,
           g.sellday,
           COALESCE(g.median::text, '無') AS median,
           COALESCE(g.count2::text, '無') AS count2,
           COALESCE(g.total_play_time_median::text, '無') AS total_play_time_median,
           COALESCE(g.time_before_understanding_fun_median::text, '無') AS time_before_understanding_fun_median,
           COALESCE(g.okazu::text, '未收錄') AS okazu,
           COALESCE(g.erogame::text, '未收錄') AS erogame,
           COALESCE(g.banner_url, '') AS banner_url,
           COALESCE(g.genre, '無') AS genre,
           COALESCE(g.steam::text, '') AS steam,
           COALESCE(g.vndb, '') AS vndb,
           g.shoukai,
           s.shubetu_detail
    FROM gamelist g
    LEFT JOIN shokushu_agg s ON s.game_id = g.id
    LEFT JOIN brandlist b ON b.id = g.brandname
    WHERE g.gamename ILIKE '%s'
    ORDER BY g.count2 DESC NULLS LAST, g.median DESC NULLS LAST
    LIMIT 1
) t;
`, result), nil
}
