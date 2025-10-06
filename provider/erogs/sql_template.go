package erogs

import (
	"fmt"
	"strings"

	kurohelpererrors "kurohelper/errors"
	"kurohelper/utils"
)

func buildSearchStringSQL(search string) (string, error) {
	search = strings.ReplaceAll(search, "'", "''")
	if strings.TrimSpace(search) == "" {
		return "", kurohelpererrors.ErrSearchNoContent
	}

	result := "%"
	searchRune := []rune(search)
	for i, r := range searchRune {
		if utils.IsEnglish(r) && i < len(searchRune)-1 {
			if utils.IsEnglish(searchRune[i+1]) {
				result += string(r)
			} else {
				result += string(r) + "%"
			}
		} else {
			result += string(r) + "%"
		}
	}
	return result, nil
}

func buildFuzzySearchCreatorSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
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
                    g.count2,
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
                GROUP BY g.id, g.gamename, g.sellday, g.median, g.count2
            ) AS game_data
        ) AS games
    FROM createrlist cr
    WHERE cr.name ILIKE '%s' OR cr.name ILIKE '%s'
    LIMIT 1
) AS c;`, resultTW, resultJP), nil
}

func buildFuzzySearchMusicSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
WITH filtered_music AS (
    SELECT 
        m.id AS music_id,
        m.name AS musicname,
        m.playtime,
        m.releasedate,
        ROUND(AVG(LEAST(ut.tokuten, 100))::numeric, 2) AS avg_tokuten,
        COUNT(DISTINCT ut.uid) AS tokuten_count
    FROM musiclist m
    LEFT JOIN usermusic_tokuten ut ON ut.music = m.id
    WHERE m.name ILIKE '%s' OR m.name ILIKE '%s'
    GROUP BY m.id, m.name, m.playtime, m.releasedate
    ORDER BY tokuten_count DESC NULLS LAST, avg_tokuten DESC NULLS LAST
    LIMIT 1
)
SELECT row_to_json(t)
FROM (
    SELECT 
        m.music_id,
        m.musicname,
        m.playtime,
        m.releasedate,
        m.avg_tokuten,
        m.tokuten_count,
        COALESCE(STRING_AGG(DISTINCT s_c.name, ','), '無') AS singer_name,
        COALESCE(STRING_AGG(DISTINCT l_c.name, ','), '無') AS lyric_name,
        COALESCE(STRING_AGG(DISTINCT a_c.name, ','), '無') AS arrangement_name,
        COALESCE(STRING_AGG(DISTINCT comp_c.name, ','), '無') AS composition_name,
        json_agg(
            DISTINCT jsonb_build_object(
                'game_name', g.gamename,
                'game_model', g.model,
                'category', gm.category
            )
        ) AS game_categories,
        COALESCE(STRING_AGG(DISTINCT mi.name, ','), '') AS album_name
    FROM filtered_music m
    LEFT JOIN singer s ON s.music = m.music_id
    LEFT JOIN createrlist s_c ON s_c.id = s.creater
    LEFT JOIN lyrics l ON l.music = m.music_id
    LEFT JOIN createrlist l_c ON l_c.id = l.creater
    LEFT JOIN arrangement a ON a.music = m.music_id
    LEFT JOIN createrlist a_c ON a_c.id = a.creater
    LEFT JOIN composition comp ON comp.music = m.music_id
    LEFT JOIN createrlist comp_c ON comp_c.id = comp.creater
    LEFT JOIN game_music gm ON gm.music = m.music_id
    LEFT JOIN gamelist g ON g.id = gm.game
    LEFT JOIN musicitem_music mim ON mim.music = m.music_id
    LEFT JOIN musicitemlist mi ON mi.id = mim.musicitem
    GROUP BY m.music_id, m.musicname, m.playtime, m.releasedate, m.avg_tokuten, m.tokuten_count
    ORDER BY tokuten_count DESC NULLS LAST, avg_tokuten DESC NULLS LAST
) t;
`, resultTW, resultJP), nil
}

func buildFuzzySearchMusicListSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
WITH filtered_music AS (
    SELECT 
        m.id AS music_id,
        m.name AS musicname,
        ROUND(AVG(LEAST(ut.tokuten, 100))::numeric, 2) AS avg_tokuten,
        COUNT(DISTINCT ut.uid) AS tokuten_count
    FROM musiclist m
    LEFT JOIN usermusic_tokuten ut ON ut.music = m.id
    WHERE m.name ILIKE '%s' OR m.name ILIKE '%s'
    GROUP BY m.id, m.name, m.playtime, m.releasedate
    ORDER BY tokuten_count DESC NULLS LAST, avg_tokuten DESC NULLS LAST
    LIMIT 200
)
SELECT json_agg(row_to_json(t))
FROM (
    SELECT 
        m.music_id AS id,
        m.musicname AS name,
        m.tokuten_count,
        m.avg_tokuten,
        STRING_AGG(DISTINCT gm.category, ',') AS category
    FROM filtered_music m
    LEFT JOIN game_music gm ON gm.music = m.music_id
    GROUP BY m.music_id,m.musicname,m.tokuten_count,m.avg_tokuten
    ORDER BY tokuten_count DESC NULLS LAST, avg_tokuten DESC NULLS LAST
) t;
`, resultTW, resultJP), nil
}

func buildFuzzySearchGameSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
WITH filtered_games AS (
    SELECT *
    FROM gamelist
    WHERE gamename ILIKE '%s' OR gamename ILIKE '%s'
    ORDER BY count2 DESC NULLS LAST, median DESC NULLS LAST
    LIMIT 1
)
SELECT row_to_json(t)
FROM (
    SELECT g.id,
           b.brandname,
           g.gamename,
           g.sellday,
           g.model AS category,
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
           j.junni,
           g.shoukai,
           s.shubetu_detail
    FROM filtered_games g
    LEFT JOIN LATERAL (
        SELECT json_agg(
                   json_build_object(
                       'shubetu_type', s.shubetu,
                       'creater_name', c.name,
                       'shubetu_detail_type', s.shubetu_detail,
                       'shubetu_detail_name', s.shubetu_detail_name
                   )
               ) AS shubetu_detail
        FROM shokushu s
        LEFT JOIN createrlist c ON c.id = s.creater
        WHERE s.game = g.id AND s.shubetu != 7
    ) s ON TRUE
    LEFT JOIN brandlist b ON b.id = g.brandname
    LEFT JOIN LATERAL (
        SELECT j.junni
        FROM junirirekimedian j
        WHERE j.game = g.id
        ORDER BY j.date DESC NULLS LAST
        LIMIT 1
    ) j ON TRUE
) t;
`, resultTW, resultJP), nil
}

func buildFuzzySearchGameListSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
SELECT json_agg(row_to_json(t))
FROM (
    SELECT g.id,
           g.gamename AS name,
           g.model AS category
    FROM gamelist g
    WHERE gamename ILIKE '%s' OR gamename ILIKE '%s'
    ORDER BY count2 DESC NULLS LAST, median DESC NULLS LAST
    LIMIT 200
) t;
`, resultTW, resultJP), nil
}

func buildFuzzySearchBrandSQL(searchTW string, searchJP string) (string, error) {
	resultTW, err := buildSearchStringSQL(searchTW)
	if err != nil {
		return "", err
	}

	resultJP, err := buildSearchStringSQL(searchJP)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(`
WITH single_brand AS (
    SELECT
        id,
        brandname,
        brandfurigana,
        url,
        kind,
        lost,
        directlink,
        median,
        twitter,
        count2,
        count_all,
        average2,
        stdev
    FROM brandlist
    WHERE brandname ILIKE '%s' OR brandname ILIKE '%s'
    ORDER BY count2 DESC NULLS LAST, median DESC NULLS LAST
    LIMIT 1
)
SELECT row_to_json(r)
FROM (
    SELECT 
        A.id, 
        A.brandname, 
        A.brandfurigana, 
        A.url, 
        A.kind, 
        A.lost, 
        A.directlink, 
        A.median, 
        A.twitter, 
        A.count2, 
        A.count_all, 
        A.average2, 
        A.stdev,
        (
            SELECT json_agg(
                json_build_object(
                    'id', g.id,
                    'gamename', g.gamename,
                    'furigana', g.furigana,
                    'sellday', g.sellday,
                    'median', g.median,
                    'model', g.model,
                    'stdev', g.stdev,
                    'count2', g.count2,
                    'vndb', g.vndb
                ) ORDER BY g.sellday DESC
            )
            FROM gamelist g
            WHERE g.brandname = A.id
        ) AS gamelist
    FROM single_brand A
) r;
`, resultTW, resultJP), nil
}
