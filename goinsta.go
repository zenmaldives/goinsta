// goinsta project goinsta.go
package goinsta

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"time"

	"github.com/erikdubbelboer/fasthttp"
)

// GetSessions return current instagram session and cookies
// Maybe need for webpages that use this API
// TODO: Adapt to goinsta...
func (insta *Instagram) GetSessions() map[string]*fasthttp.Cookie {
	return insta.cookies.Cookies()
}

// SetCookies can enable us to set cookie, it'll be help for webpage that use this API without Login-again.
func (insta *Instagram) SetCookies(cks []*fasthttp.Cookie) {
	if insta.cookies == nil {
		insta.cookies = &cookies{}
	}
	insta.cookies.SetCookies(cks)
}

// Const values ,
// GOINSTA Default variables contains API url , user agent and etc...
// goInstaSigKey is Instagram sign key, It's important
// Filter_<name>
const (
	Filter_Walden        = 20
	Filter_Crema         = 616
	Filter_Reyes         = 614
	Filter_Moon          = 111
	Filter_Ashby         = 116
	Filter_Maven         = 118
	Filter_Brannan       = 22
	Filter_Hefe          = 21
	Filter_Valencia      = 25
	Filter_Clarendon     = 112
	Filter_Helena        = 117
	Filter_Brooklyn      = 115
	Filter_Dogpatch      = 105
	Filter_Ludwig        = 603
	Filter_Stinson       = 109
	Filter_Inkwell       = 10
	Filter_Rise          = 23
	Filter_Perpetua      = 608
	Filter_Juno          = 613
	Filter_Charmes       = 108
	Filter_Ginza         = 107
	Filter_Hudson        = 26
	Filter_Normat        = 0
	Filter_Slumber       = 605
	Filter_Lark          = 615
	Filter_Skyline       = 113
	Filter_Kelvin        = 16
	Filter_1977          = 14
	Filter_Lo_Fi         = 2
	Filter_Aden          = 612
	Filter_Amaro         = 24
	Filter_Sutro         = 18
	Filter_Vasper        = 106
	Filter_Nashville     = 15
	Filter_X_Pro_II      = 1
	Filter_Mayfair       = 17
	Filter_Toaster       = 19
	Filter_Earlybird     = 3
	Filter_Willow        = 28
	Filter_Sierra        = 27
	Filter_Gingham       = 114
	goInstaAPIUrl        = "https://i.instagram.com/api/v1/"
	goInstaUserAgent     = "Instagram 10.26.0 Android (18/4.3; 320dpi; 720x1280; Xiaomi; HM 1SW; armani; qcom; en_US)"
	goInstaSigKey        = "4f8732eb9ba7d1c8e8897a75d6474d4eb3f5279137431b2aafb71fafe2abe178"
	GOINSTA_EXPERIMENTS  = "ig_promote_reach_objective_fix_universe,ig_android_universe_video_production,ig_search_client_h1_2017_holdout,ig_android_live_follow_from_comments_universe,ig_android_carousel_non_square_creation,ig_android_live_analytics,ig_android_follow_all_dialog_confirmation_copy,ig_android_stories_server_coverframe,ig_android_video_captions_universe,ig_android_offline_location_feed,ig_android_direct_inbox_retry_seen_state,ig_android_ontact_invite_universe,ig_android_live_broadcast_blacklist,ig_android_insta_video_reconnect_viewers,ig_android_ad_async_ads_universe,ig_android_search_clear_layout_universe,ig_android_shopping_reporting,ig_android_stories_surface_universe,ig_android_verified_comments_universe,ig_android_preload_media_ahead_in_current_reel,android_instagram_prefetch_suggestions_universe,ig_android_reel_viewer_fetch_missing_reels_universe,ig_android_direct_search_share_sheet_universe,ig_android_business_promote_tooltip,ig_android_direct_blue_tab,ig_android_async_network_tweak_universe,ig_android_elevate_main_thread_priority_universe,ig_android_stories_gallery_nux,ig_android_instavideo_remove_nux_comments,ig_video_copyright_whitelist,ig_react_native_inline_insights_with_relay,ig_android_direct_thread_message_animation,ig_android_draw_rainbow_client_universe,ig_android_direct_link_style,ig_android_live_heart_enhancements_universe,ig_android_rtc_reshare,ig_android_preload_item_count_in_reel_viewer_buffer,ig_android_users_bootstrap_service,ig_android_auto_retry_post_mode,ig_android_shopping,ig_android_main_feed_seen_state_dont_send_info_on_tail_load,ig_fbns_preload_default,ig_android_gesture_dismiss_reel_viewer,ig_android_tool_tip,ig_android_ad_logger_funnel_logging_universe,ig_android_gallery_grid_column_count_universe,ig_android_business_new_ads_payment_universe,ig_android_direct_links,ig_android_audience_control,ig_android_live_encore_consumption_settings_universe,ig_perf_android_holdout,ig_android_cache_contact_import_list,ig_android_links_receivers,ig_android_ad_impression_backtest,ig_android_list_redesign,ig_android_stories_separate_overlay_creation,ig_android_stop_video_recording_fix_universe,ig_android_render_video_segmentation,ig_android_live_encore_reel_chaining_universe,ig_android_sync_on_background_enhanced_10_25,ig_android_immersive_viewer,ig_android_mqtt_skywalker,ig_fbns_push,ig_android_ad_watchmore_overlay_universe,ig_android_react_native_universe,ig_android_profile_tabs_redesign_universe,ig_android_live_consumption_abr,ig_android_story_viewer_social_context,ig_android_hide_post_in_feed,ig_android_video_loopcount_int,ig_android_enable_main_feed_reel_tray_preloading,ig_android_camera_upsell_dialog,ig_android_ad_watchbrowse_universe,ig_android_internal_research_settings,ig_android_search_people_tag_universe,ig_android_react_native_ota,ig_android_enable_concurrent_request,ig_android_react_native_stories_grid_view,ig_android_business_stories_inline_insights,ig_android_log_mediacodec_info,ig_android_direct_expiring_media_loading_errors,ig_video_use_sve_universe,ig_android_cold_start_feed_request,ig_android_enable_zero_rating,ig_android_reverse_audio,ig_android_branded_content_three_line_ui_universe,ig_android_live_encore_production_universe,ig_stories_music_sticker,ig_android_stories_teach_gallery_location,ig_android_http_stack_experiment_2017,ig_android_stories_device_tilt,ig_android_pending_request_search_bar,ig_android_fb_topsearch_sgp_fork_request,ig_android_seen_state_with_view_info,ig_android_animation_perf_reporter_timeout,ig_android_new_block_flow,ig_android_story_tray_title_play_all_v2,ig_android_direct_address_links,ig_android_stories_archive_universe,ig_android_save_collections_cover_photo,ig_android_live_webrtc_livewith_production,ig_android_sign_video_url,ig_android_stories_video_prefetch_kb,ig_android_stories_create_flow_favorites_tooltip,ig_android_live_stop_broadcast_on_404,ig_android_live_viewer_invite_universe,ig_android_promotion_feedback_channel,ig_android_render_iframe_interval,ig_android_accessibility_logging_universe,ig_android_camera_shortcut_universe,ig_android_use_one_cookie_store_per_user_override,ig_profile_holdout_2017_universe,ig_android_stories_server_brushes,ig_android_ad_media_url_logging_universe,ig_android_shopping_tag_nux_text_universe,ig_android_comments_single_reply_universe,ig_android_stories_video_loading_spinner_improvements,ig_android_collections_cache,ig_android_comment_api_spam_universe,ig_android_facebook_twitter_profile_photos,ig_android_shopping_tag_creation_universe,ig_story_camera_reverse_video_experiment,ig_android_direct_bump_selected_recipients,ig_android_ad_cta_haptic_feedback_universe,ig_android_vertical_share_sheet_experiment,ig_android_family_bridge_share,ig_android_search,ig_android_insta_video_consumption_titles,ig_android_stories_gallery_preview_button,ig_android_fb_auth_education,ig_android_camera_universe,ig_android_me_only_universe,ig_android_instavideo_audio_only_mode,ig_android_user_profile_chaining_icon,ig_android_live_video_reactions_consumption_universe,ig_android_stories_hashtag_text,ig_android_post_live_badge_universe,ig_android_swipe_fragment_container,ig_android_search_users_universe,ig_android_live_save_to_camera_roll_universe,ig_creation_growth_holdout,ig_android_sticker_region_tracking,ig_android_unified_inbox,ig_android_live_new_watch_time,ig_android_offline_main_feed_10_11,ig_import_biz_contact_to_page,ig_android_live_encore_consumption_universe,ig_android_experimental_filters,ig_android_search_client_matching_2,ig_android_react_native_inline_insights_v2,ig_android_business_conversion_value_prop_v2,ig_android_redirect_to_low_latency_universe,ig_android_ad_show_new_awr_universe,ig_family_bridges_holdout_universe,ig_android_background_explore_fetch,ig_android_following_follower_social_context,ig_android_video_keep_screen_on,ig_android_ad_leadgen_relay_modern,ig_android_profile_photo_as_media,ig_android_insta_video_consumption_infra,ig_android_ad_watchlead_universe,ig_android_direct_prefetch_direct_story_json,ig_android_shopping_react_native,ig_android_top_live_profile_pics_universe,ig_android_direct_phone_number_links,ig_android_stories_weblink_creation,ig_android_direct_search_new_thread_universe,ig_android_histogram_reporter,ig_android_direct_on_profile_universe,ig_android_network_cancellation,ig_android_background_reel_fetch,ig_android_react_native_insights,ig_android_insta_video_audio_encoder,ig_android_family_bridge_bookmarks,ig_android_data_usage_network_layer,ig_android_universal_instagram_deep_links,ig_android_dash_for_vod_universe,ig_android_modular_tab_discover_people_redesign,ig_android_mas_sticker_upsell_dialog_universe,ig_android_ad_add_per_event_counter_to_logging_event,ig_android_sticky_header_top_chrome_optimization,ig_android_rtl,ig_android_biz_conversion_page_pre_select,ig_android_promote_from_profile_button,ig_android_live_broadcaster_invite_universe,ig_android_share_spinner,ig_android_text_action,ig_android_own_reel_title_universe,ig_promotions_unit_in_insights_landing_page,ig_android_business_settings_header_univ,ig_android_save_longpress_tooltip,ig_android_constrain_image_size_universe,ig_android_business_new_graphql_endpoint_universe,ig_ranking_following,ig_android_stories_profile_camera_entry_point,ig_android_universe_reel_video_production,ig_android_power_metrics,ig_android_sfplt,ig_android_offline_hashtag_feed,ig_android_live_skin_smooth,ig_android_direct_inbox_search,ig_android_stories_posting_offline_ui,ig_android_sidecar_video_upload_universe,ig_android_promotion_manager_entry_point_universe,ig_android_direct_reply_audience_upgrade,ig_android_swipe_navigation_x_angle_universe,ig_android_offline_mode_holdout,ig_android_live_send_user_location,ig_android_direct_fetch_before_push_notif,ig_android_non_square_first,ig_android_insta_video_drawing,ig_android_swipeablefilters_universe,ig_android_live_notification_control_universe,ig_android_analytics_logger_running_background_universe,ig_android_save_all,ig_android_reel_viewer_data_buffer_size,ig_direct_quality_holdout_universe,ig_android_family_bridge_discover,ig_android_react_native_restart_after_error_universe,ig_android_startup_manager,ig_story_tray_peek_content_universe,ig_android_profile,ig_android_high_res_upload_2,ig_android_http_service_same_thread,ig_android_scroll_to_dismiss_keyboard,ig_android_remove_followers_universe,ig_android_skip_video_render,ig_android_story_timestamps,ig_android_live_viewer_comment_prompt_universe,ig_profile_holdout_universe,ig_android_react_native_insights_grid_view,ig_stories_selfie_sticker,ig_android_stories_reply_composer_redesign,ig_android_streamline_page_creation,ig_explore_netego,ig_android_ig4b_connect_fb_button_universe,ig_android_feed_util_rect_optimization,ig_android_rendering_controls,ig_android_os_version_blocking,ig_android_encoder_width_safe_multiple_16,ig_search_new_bootstrap_holdout_universe,ig_android_snippets_profile_nux,ig_android_e2e_optimization_universe,ig_android_comments_logging_universe,ig_shopping_insights,ig_android_save_collections,ig_android_live_see_fewer_videos_like_this_universe,ig_android_show_new_contact_import_dialog,ig_android_live_view_profile_from_comments_universe,ig_fbns_blocked,ig_formats_and_feedbacks_holdout_universe,ig_android_reduce_view_pager_buffer,ig_android_instavideo_periodic_notif,ig_search_user_auto_complete_cache_sync_ttl,ig_android_marauder_update_frequency,ig_android_suggest_password_reset_on_oneclick_login,ig_android_promotion_entry_from_ads_manager_universe,ig_android_live_special_codec_size_list,ig_android_enable_share_to_messenger,ig_android_background_main_feed_fetch,ig_android_live_video_reactions_creation_universe,ig_android_channels_home,ig_android_sidecar_gallery_universe,ig_android_upload_reliability_universe,ig_migrate_mediav2_universe,ig_android_insta_video_broadcaster_infra_perf,ig_android_business_conversion_social_context,android_ig_fbns_kill_switch,ig_android_live_webrtc_livewith_consumption,ig_android_destroy_swipe_fragment,ig_android_react_native_universe_kill_switch,ig_android_stories_book_universe,ig_android_all_videoplayback_persisting_sound,ig_android_draw_eraser_universe,ig_direct_search_new_bootstrap_holdout_universe,ig_android_cache_layer_bytes_threshold,ig_android_search_hash_tag_and_username_universe,ig_android_business_promotion,ig_android_direct_search_recipients_controller_universe,ig_android_ad_show_full_name_universe,ig_android_anrwatchdog,ig_android_qp_kill_switch,ig_android_2fac,ig_direct_bypass_group_size_limit_universe,ig_android_promote_simplified_flow,ig_android_share_to_whatsapp,ig_android_hide_bottom_nav_bar_on_discover_people,ig_fbns_dump_ids,ig_android_hands_free_before_reverse,ig_android_skywalker_live_event_start_end,ig_android_live_join_comment_ui_change,ig_android_direct_search_story_recipients_universe,ig_android_direct_full_size_gallery_upload,ig_android_ad_browser_gesture_control,ig_channel_server_experiments,ig_android_video_cover_frame_from_original_as_fallback,ig_android_ad_watchinstall_universe,ig_android_ad_viewability_logging_universe,ig_android_new_optic,ig_android_direct_visual_replies,ig_android_stories_search_reel_mentions_universe,ig_android_threaded_comments_universe,ig_android_mark_reel_seen_on_Swipe_forward,ig_internal_ui_for_lazy_loaded_modules_experiment,ig_fbns_shared,ig_android_capture_slowmo_mode,ig_android_live_viewers_list_search_bar,ig_android_video_single_surface,ig_android_offline_reel_feed,ig_android_video_download_logging,ig_android_last_edits,ig_android_exoplayer_4142,ig_android_post_live_viewer_count_privacy_universe,ig_android_activity_feed_click_state,ig_android_snippets_haptic_feedback,ig_android_gl_drawing_marks_after_undo_backing,ig_android_mark_seen_state_on_viewed_impression,ig_android_live_backgrounded_reminder_universe,ig_android_live_hide_viewer_nux_universe,ig_android_live_monotonic_pts,ig_android_search_top_search_surface_universe,ig_android_user_detail_endpoint,ig_android_location_media_count_exp_ig,ig_android_comment_tweaks_universe,ig_android_ad_watchmore_entry_point_universe,ig_android_top_live_notification_universe,ig_android_add_to_last_post,ig_save_insights,ig_android_live_enhanced_end_screen_universe,ig_android_ad_add_counter_to_logging_event,ig_android_blue_token_conversion_universe,ig_android_exoplayer_settings,ig_android_progressive_jpeg,ig_android_offline_story_stickers,ig_android_gqls_typing_indicator,ig_android_chaining_button_tooltip,ig_android_video_prefetch_for_connectivity_type,ig_android_use_exo_cache_for_progressive,ig_android_samsung_app_badging,ig_android_ad_holdout_watchandmore_universe,ig_android_offline_commenting,ig_direct_stories_recipient_picker_button,ig_insights_feedback_channel_universe,ig_android_insta_video_abr_resize,ig_android_insta_video_sound_always_on"
	goInstaSigKeyVersion = "4"
)

// GOINSTA_DEVICE_SETTINGS variable is a simulate of an android device
var GOINSTA_DEVICE_SETTINGS = map[string]interface{}{
	"manufacturer":    "Xiaomi",
	"model":           "HM 1SW",
	"android_version": 18,
	"android_release": "4.3",
}

// NewViaProxy All requests will use proxy server (example http://<ip>:<port>)
func NewViaProxy(username, password string, dialFunc fasthttp.DialFunc) *Instagram {
	insta := New(username, password)
	insta.DialFunc = dialFunc
	return insta
}

// New try to fill Instagram struct
// New does not try to login , it will only fill
// Instagram struct
func New(username, password string) *Instagram {
	insta := &Instagram{
		client: &fasthttp.Client{
			Name: goInstaUserAgent,
		},
		cookies: nil,
		Info: &ClientInfo{
			DeviceID: generateDeviceID(generateMD5Hash(username + password)),
			Username: username,
			Password: password,
			UUID:     generateUUID(true),
			PhoneID:  generateUUID(true),
		},
	}
	insta.fill()
	return insta
}

// Login to Instagram.
// return error if can't send request to instagram server
func (insta *Instagram) Login() error {
	req := acquireRequest()
	defer releaseRequest(req)

	req.SetEndpoint("si/fetch_headers/")
	req.args.Set("challenge_type", "signup")
	req.args.Set("guid", generateUUID(false))

	body, err := insta.sendRequest(req)
	if err != nil {
		return fmt.Errorf("login failed for %s: %s", insta.Info.Username, err)
	}

	result, _ := json.Marshal(map[string]interface{}{
		"guid":                insta.Info.UUID,
		"login_attempt_count": 0,
		"_csrftoken":          insta.Info.Token,
		"device_id":           insta.Info.DeviceID,
		"phone_id":            insta.Info.PhoneID,
		"username":            insta.Info.Username,
		"password":            insta.Info.Password,
	})

	req.SetEndpoint("accounts/login/")
	req.SetData(generateSignature(b2s(result)))
	body, err = insta.sendRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, insta)
	if err != nil {
		return err
	}

	//insta.Current = &Result.Current
	insta.Logged = true
	insta.Info.RankToken = strconv.FormatInt(Result.Current.ID, 10) + "_" + insta.Info.UUID

	insta.SyncFeatures()
	insta.AutoCompleteUserList()
	insta.GetRankedRecipients()
	insta.Timeline("")
	insta.GetRankedRecipients()
	insta.GetRecentRecipients()
	insta.MegaphoneLog()
	insta.GetV2Inbox()
	insta.GetRecentActivity()
	insta.GetReelsTrayFeed()

	return nil
}

// Logout of Instagram
func (insta *Instagram) Logout() error {
	_, err := insta.sendSimpleRequest("accounts/logout/")
	if err == nil {
		insta.Logged = false
		insta.Info = nil
		insta.Current = nil
		insta.client = nil
		insta.cookies.Release()
		insta.cookies = nil
	}
	return err
}

// SyncFeatures simulates Instagram app behavior
func (insta *Instagram) SyncFeatures() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"id":          insta.Current.ID,
			"experiments": GOINSTA_EXPERIMENTS,
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)

	req.SetEndpoint("qe/sync/")
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// AutoCompleteUserList simulates Instagram app behavior
func (insta *Instagram) AutoCompleteUserList() error {
	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)

	req.SetEndpoint("friendships/autocomplete_user_list/")
	req.skipStatus = true
	req.args.Set("version", "2")
	_, err := insta.sendRequest(req)
	return err
}

// MegaphoneLog simulates Instagram app behavior
func (insta *Instagram) MegaphoneLog() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"id":        insta.Current.ID,
			"type":      "feed_aysf",
			"action":    "seen",
			"reason":    "",
			"device_id": insta.Info.DeviceID,
			"uuid":      generateMD5Hash(string(time.Now().Unix())),
		},
	)
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)

	req.SetEndpoint("megaphone/log/")
	req.SetData(generateSignature(data))

	_, err = insta.sendRequest(req)
	return err
}

// Expose , expose instagram
// return error if status was not 'ok' or runtime error
func (insta *Instagram) Expose() error {
	data, err := insta.prepareData(
		map[string]interface{}{
			"id":         insta.Current.ID,
			"experiment": "ig_android_profile_contextual_feed",
		},
	)
	if err != nil {
		return err
	}

	result := StatusResponse{}

	req := acquireRequest()
	defer releaseRequest(req)

	req.SetEndpoint("qe/expose/")
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &result)
}

// TODO
func (insta *Instagram) fill() {
	if insta.User == nil {
		insta.User = NewUser(insta)
	}
	user := insta.User
	if insta.Current == nil {
		insta.Current = &ProfileData{}
	}

	if insta.Tag == nil {
		insta.Tag = NewTag(insta)
	}
	if insta.Current.Feed == nil {
		insta.Current.Feed = NewFeed(user)
	}
	if insta.Inbox == nil {
		insta.Inbox = NewInbox(insta)
	}
	if insta.Current.Following == nil {
		insta.Current.Following = NewUsers(user, false)
	}
	if insta.Current.Followers == nil {
		insta.Current.Followers = NewUsers(user, true)
	}

	if insta.Current.insta == nil {
		insta.Current.insta = insta
	}
	if insta.Media == nil {
		insta.Media = NewMedia(insta)
	}
	if insta.Search == nil {
		insta.Search = NewSearch(insta)
	}
}

// SetPublicAccount sets account to public
func (insta *Instagram) SetPublicAccount() error {
	defer insta.fill()

	data, err := insta.prepareData(make(map[string]interface{}))
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint("accounts/set_public/")
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, insta)
}

// SetPrivateAccount sets account to private
func (insta *Instagram) SetPrivateAccount() error {
	defer insta.fill()

	data, err := insta.prepareData(make(map[string]interface{}))
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint("accounts/set_private/")
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, insta)
}

// GetProfileData return current user information
func (insta *Instagram) GetProfileData() error {
	defer insta.fill()

	data, err := insta.prepareData(make(map[string]interface{}))
	if err != nil {
		return err
	}

	req := acquireRequest()
	req.args = fasthttp.AcquireArgs()
	defer releaseRequest(req)
	req.SetEndpoint("accounts/current_user/")
	req.SetData(generateSignature(data))
	req.args.Set("edit", "true")

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, insta)
}

// RemoveProfilePicture will remove current logged in user profile picture
func (insta *Instagram) RemoveProfilePicture() error {
	data, err := insta.prepareData(make(map[string]interface{}))
	if err != nil {
		return err
	}

	req := acquireRequest()
	defer releaseRequest(req)
	req.SetEndpoint("accounts/remove_profile_picture/")
	req.SetData(generateSignature(data))

	body, err := insta.sendRequest(req)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, insta)
}

// getImageDimensionFromReader return image dimension , types is .jpg and .png
func getImageDimensionFromReader(rdr io.Reader) (int, int, error) {
	image, _, err := image.DecodeConfig(rdr)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
}

// getImageDimension return image dimension , types is .jpg and .png
func getImageDimension(imagePath string) (int, int, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	return getImageDimensionFromReader(file)
}

// Explore stores explore menu in Instagram.Explore
func (insta *Instagram) Explore() error {
	body, err := insta.sendSimpleRequest("discover/explore/")
	if err != nil {
		return err
	}
	return json.Unmarshal(body, &insta.Explore)
}

func (insta *Instagram) prepareData(data map[string]interface{}) (string, error) {
	data["_uuid"] = insta.Info.UUID
	data["_uid"] = insta.CurrentUser.ID
	data["_csrftoken"] = insta.Info.Token
	bytes, err := json.Marshal(data)
	return b2s(bytes), err
}
