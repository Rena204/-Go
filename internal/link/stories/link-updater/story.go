package link_updater

import (
	"context"
)

func New(repository repository, consumer amqpConsumer) *Story {
	return &Story{repository: repository, consumer: consumer}
}

type Story struct {
	repository repository
	consumer   amqpConsumer
}

func (s *Story) Run(ctx context.Context) error {
	    // Listen to the queue
		msgs, err := s.consumer.Consume()
		if err != nil {
			return err
		}
	
		for msg := range msgs {
			// Decode the message
			var m message
			err := json.Unmarshal(msg.Body, &m)
			if err != nil {
				return err
			}
	
			// Get the link object from the repository
			link, err := s.repository.GetLink(m.ID)
			if err != nil {
				return err
			}
	
			// Call the scrape package
			scrapedData, err := scrape.Scrape(link.URL)
			if err != nil {
				return err
			}
	
			// If new data is parsed, add it to the link object
			if scrapedData.Title != "" {
				link.Title = scrapedData.Title
			}
			if scrapedData.Description != "" {
				link.Description = scrapedData.Description
			}
	
			// Update the link object in the repository
			err = s.repository.UpdateLink(link)
			if err != nil {
				return err
			}
		}
	
		return nil
	}
